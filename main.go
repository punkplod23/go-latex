package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func generatePDFHandler(w http.ResponseWriter, r *http.Request) {
	// LaTeX content including TikZ/pictures support for graphs
	texContent := `
	\documentclass{article}
	\usepackage{tikz}
	\usepackage{pgfplots}
	\usepackage{pgf-pie}
	\pgfplotsset{compat=1.18}
	
	\begin{document}
	\section*{Streamed PDF with Charts}
	This document and its charts were compiled in memory inside the Ubuntu container.

	\begin{center}
	\begin{tikzpicture}
	\begin{axis}[
	    ybar,
	    enlargelimits=0.15,
	    legend style={at={(0.5,-0.15)},
	      anchor=north,legend columns=-1},
	    ylabel={Performance Metric},
	    symbolic x coords={Q1,Q2,Q3,Q4},
	    xtick=data,
	    nodes near coords,
	    nodes near coords align={vertical},
	    ymin=0,
	]
	\addplot coordinates {(Q1,25) (Q2,45) (Q3,60) (Q4,85)};
	\addplot coordinates {(Q1,20) (Q2,35) (Q3,50) (Q4,75)};
	\legend{Legacy, Current}
	\end{axis}
	\end{tikzpicture}
	\end{center}

	\newpage
	\section*{Pie Chart}
	\begin{center}
	\begin{tikzpicture}
	\pie[sum=auto, text=legend, radius=2.5]{40/Legacy, 35/Current, 25/Other}
	\end{tikzpicture}
	\end{center}

	\section*{Horizontal Bar Chart}
	\begin{center}
	\begin{tikzpicture}
	\begin{axis}[
	    xbar,
	    width=11cm,
	    height=6cm,
	    enlarge y limits=0.2,
	    xlabel={Performance Metric},
	    symbolic y coords={Infrastructure, Documentation, Testing, Deployment},
	    ytick=data,
	    nodes near coords,
	    xmin=0,
	]
	\addplot coordinates {
	    (65,Infrastructure)
	    (80,Documentation)
	    (55,Testing)
	    (90,Deployment)
	};
	\end{axis}
	\end{tikzpicture}
	\end{center}
	\end{document}
	`

	// pdflatex writes the PDF to a file, not to stdout. Compile in an isolated
	// temporary directory so generated auxiliary files do not pollute the app.
	tempDir, err := os.MkdirTemp("", "go-latex-")
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to create temporary directory: %v", err), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)

	texPath := tempDir + "/document.tex"
	if err := os.WriteFile(texPath, []byte(texContent), 0600); err != nil {
		http.Error(w, fmt.Sprintf("Unable to write LaTeX source: %v", err), http.StatusInternalServerError)
		return
	}

	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-halt-on-error", "document.tex")
	cmd.Dir = tempDir

	// Capture both streams because pdflatex normally writes its diagnostics to stdout.
	var compilerLog bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &compilerLog
	cmd.Stderr = &stderr

	// Run compilation
	if err := cmd.Run(); err != nil {
		http.Error(w, fmt.Sprintf("Compilation failed: %v\nLogs: %s%s", err, compilerLog.String(), stderr.String()), http.StatusInternalServerError)
		log.Printf("LaTeX error: %s%s", compilerLog.String(), stderr.String())
		return
	}

	pdfBuf, err := os.ReadFile(tempDir + "/document.pdf")
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read generated PDF: %v\nLogs: %s%s", err, compilerLog.String(), stderr.String()), http.StatusInternalServerError)
		return
	}

	// Set headers for PDF streaming
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=generated.pdf")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBuf)))

	// Stream the compiled PDF directly to the response writer
	if _, err := io.Copy(w, bytes.NewReader(pdfBuf)); err != nil {
		log.Printf("Failed to stream PDF: %v", err)
	}
}

func main() {
	http.HandleFunc("/generate", generatePDFHandler)
	log.Println("Server listening on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
