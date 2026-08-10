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
	\documentclass[a4paper]{article}
    \usepackage[top=3cm, bottom=3cm, left=1cm, right=1cm]{geometry}
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
    \definecolor{paletteExtra}{HTML}{8B5CF6}

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
			area legend,
			axis x line=bottom,
			axis y line=left,
			tick style={draw=paletteInk},
			enlargelimits=0.15,
			legend style={draw=none, fill=none, at={(0.5,-0.22)},
			anchor=north, legend columns=4},
			ylabel={Performance Metric},
			symbolic x coords={Q1,Q2,Q3,Q4},
			xtick={Q1,Q2,Q3,Q4},
			ymin=0,
			enlarge y limits=false,
			bar width=3cm,
			bar shift=0pt,
			width=0.85\textwidth,
			height=7cm,
		]
			\addplot[fill=palettePrimary, draw=none] coordinates {(Q1,20)};
			\addlegendentry{Q1}

			\addplot[fill=paletteSecondary, draw=none] coordinates {(Q2,30)};
			\addlegendentry{Q2}

			\addplot[fill=paletteAccent, draw=none] coordinates {(Q3,40)};
			\addlegendentry{Q3}

			\addplot[fill=paletteExtra, draw=none] coordinates {(Q4,50)};
			\addlegendentry{Q4}

		\end{axis}
		\end{tikzpicture}
	\end{center}
    
	\newpage
    \section*{Pie Chart}
    \begin{center}
    \begin{tikzpicture}
    \pie[sum=auto, text=pin, hide number, line width=0pt, radius=4.0, color={palettePrimary,paletteSecondary,paletteAccent}]{40/{Legacy: 40\%}, 35/{Current: 35\%}, 25/{Other: 25\%}}
    \begin{scope}[yshift=-4.8cm]
        \fill[palettePrimary] (-2.4,0) rectangle (-2.2,0.2);
        \node[anchor=west, font=\small] at (-2.1,0.1) {Legacy};
        \fill[paletteSecondary] (0,0) rectangle (0.2,0.2);
        \node[anchor=west, font=\small] at (0.3,0.1) {Current};
        \fill[paletteAccent] (2.0,0) rectangle (2.2,0.2);
        \node[anchor=west, font=\small] at (2.3,0.1) {Other};
    \end{scope}
    \end{tikzpicture}
    \end{center}

	\newpage
    \section*{Horizontal Bar Chart}
	\begin{center}
	\begin{tikzpicture}
	\begin{axis}[
	    xbar,
		area legend,
	    axis x line=top,
	    axis y line=left,
	    y dir=reverse,
	    tick style={draw=paletteInk},
		legend style={draw=none, fill=none, at={(0.5,-0.22)}, anchor=north, legend columns=4},
	    xmajorgrids=true,
	    ymajorgrids=false,
	    grid style={draw=paletteGrid, line width=0.4pt},
	    xtick pos=top,
	    xlabel style={at={(axis description cs:0.5,1.12)}, anchor=south, text=paletteInk},
	    x tick label style={anchor=south, yshift=2pt, text=paletteInk},
		width=0.85\textwidth,
	    enlarge y limits=0.2,
        enlarge x limits=0,
	    bar width=18pt,
	    xlabel={Performance Metric},
	    symbolic y coords={Infrastructure, Documentation, Testing, Deployment},
		ytick={Infrastructure, Documentation, Testing, Deployment},
	    yticklabel style={text=paletteInk},
		xmin=0,
		height=6cm,
		xmax=100,
	]
		\addplot[fill=palettePrimary, draw=none, bar shift=0pt] coordinates {(65,Infrastructure)};
		\addlegendentry{Infrastructure}

		\addplot[fill=paletteSecondary, draw=none, bar shift=0pt] coordinates {(80,Documentation)};
		\addlegendentry{Documentation}

		\addplot[fill=paletteAccent, draw=none, bar shift=0pt] coordinates {(55,Testing)};
		\addlegendentry{Testing}

		\addplot[fill=paletteExtra, draw=none, bar shift=0pt] coordinates {(90,Deployment)};
		\addlegendentry{Deployment}
	\end{axis}
	\end{tikzpicture}
	\end{center}

    \vspace{1cm}
    \section*{Summary Table}
	\begin{center}
	\renewcommand{\arraystretch}{1.5} % Increases vertical row padding (default is 1.0)
	\begin{tabular}{|p{0.4\textwidth}|p{0.25\textwidth}|p{0.25\textwidth}|}
	\hline
	\textbf{Quarter / Metric} & \textbf{Status} & \textbf{Value} \\ \hline
	Q1 Performance & Complete & 20 \\ \hline
	Q2 Performance & Complete & 35 \\ \hline
	Q3 Performance & Complete & 50 \\ \hline
	\end{tabular}
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
