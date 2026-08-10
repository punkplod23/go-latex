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
	\usepackage{fancyhdr}
	\usepackage{fontspec}
	\setmainfont{Roboto}
	\pgfplotsset{compat=1.18}

	% Configure all report colors here.
	\definecolor{palettePrimary}{HTML}{2563EB}
	\definecolor{paletteSecondary}{HTML}{14B8A6}
	\definecolor{paletteAccent}{HTML}{F59E0B}
	\definecolor{paletteInk}{HTML}{1E293B}
	\definecolor{paletteGrid}{HTML}{CBD5E1}

	\begin{document}
	\pagestyle{fancy}
	\fancyhf{}
	\lhead{Chart Performance Report}
	\rhead{Go + LaTeX}
	\cfoot{Page \thepage}
	\setlength{\headheight}{14pt}

	\begin{titlepage}
	\centering
	\vspace*{3cm}
	{\Huge\bfseries Chart Performance Report\par}
	\vspace{1.5cm}
	{\Large Generated PDF Report\par}
	\vspace{2cm}
	\begin{tikzpicture}[
	    neuron/.style={circle, draw=palettePrimary, fill=palettePrimary!15, minimum size=9mm, thick},
	    connection/.style={draw=palettePrimary!45, thick}
	]
	\foreach \y in {1,2,3,4}
	    \node[neuron] (input\y) at (0,\y) {};
	\foreach \y in {1,2,3,4,5}
	    \node[neuron] (hidden\y) at (2.5,\y-0.5) {};
	\foreach \y in {1,2,3}
	    \node[neuron] (output\y) at (5,\y+0.5) {};
	\foreach \i in {1,2,3,4}
	    \foreach \j in {1,2,3,4,5}
	        \draw[connection] (input\i) -- (hidden\j);
	\foreach \i in {1,2,3,4,5}
	    \foreach \j in {1,2,3}
	        \draw[connection] (hidden\i) -- (output\j);
	\end{tikzpicture}
	\par\vspace{1cm}
	{\large Compiled with Go and LaTeX\par}
	\vfill
	{\large \today\par}
	\end{titlepage}

	\section*{Streamed PDF with Charts}
	This document and its charts were compiled in memory inside the Alpine container.

	\begin{center}
	\begin{tikzpicture}
	\begin{axis}[
	    ybar,
	    axis x line=bottom,
	    axis y line=left,
	    tick style={draw=paletteInk},
	    cycle list={{fill=palettePrimary, draw=none}, {fill=paletteSecondary, draw=none}},
	    enlargelimits=0.15,
	    legend style={draw=none, fill=none, at={(0.5,-0.15)},
	      anchor=north,legend columns=-1},
	    ylabel={Performance Metric},
	    symbolic x coords={Q1,Q2,Q3,Q4},
	    xtick=data,
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
	\pie[sum=auto, text=pin, hide number, line width=0pt, radius=2.5, color={palettePrimary,paletteSecondary,paletteAccent}]{40/{Legacy: 40\%}, 35/{Current: 35\%}, 25/{Other: 25\%}}
	\begin{scope}[yshift=-3.4cm]
	    \fill[palettePrimary] (-2.4,0) rectangle (-2.2,0.2);
	    \node[anchor=west, font=\small] at (-2.1,0.1) {Legacy};
	    \fill[paletteSecondary] (0,0) rectangle (0.2,0.2);
	    \node[anchor=west, font=\small] at (0.3,0.1) {Current};
	    \fill[paletteAccent] (2.0,0) rectangle (2.2,0.2);
	    \node[anchor=west, font=\small] at (2.3,0.1) {Other};
	\end{scope}
	\end{tikzpicture}
	\end{center}

	\section*{Horizontal Bar Chart}
	\begin{center}
	\begin{tikzpicture}
	\begin{axis}[
	    xbar,
	    axis x line=top,
	    axis y line=left,
	    tick style={draw=paletteInk},
	    cycle list={{fill=paletteAccent, draw=none}},
	    legend style={draw=none, fill=none, at={(0.5,-0.2)}, anchor=north, text=paletteInk},
	    xmajorgrids=true,
	    ymajorgrids=false,
	    grid style={draw=paletteGrid, line width=0.4pt},
	    xtick pos=top,
	    xlabel style={at={(axis description cs:0.5,1.12)}, anchor=south, text=paletteInk},
	    x tick label style={anchor=south, yshift=2pt, text=paletteInk},
	    width=11cm,
	    height=6cm,
	    enlarge y limits=0.2,
	    xlabel={Performance Metric},
	    symbolic y coords={Infrastructure, Documentation, Testing, Deployment},
	    ytick=data,
	    yticklabel style={text=paletteInk},
	    xmin=0,
	]
	\addplot coordinates {
	    (65,Infrastructure)
	    (80,Documentation)
	    (55,Testing)
	    (90,Deployment)
	};
	\legend{Performance Metric}
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

	cmd := exec.Command("xelatex", "-interaction=nonstopmode", "-halt-on-error", "-no-shell-escape", "document.tex")
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
