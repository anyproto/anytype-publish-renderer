function initToggles() {
    const toggles = document.querySelectorAll(".textToggle");
    toggles.forEach(t => {
        t.addEventListener("click", function() {
            t.classList.toggle("isToggled");
        })
    })

}

function initLatex() {
    const katex = window.katex
    const latexBlocks = document.querySelectorAll(".isLatex .content");
    const trustFn = context => {
        return [ '\\url', '\\href', '\\includegraphics' ].includes(context.command)
    }
    latexBlocks.forEach(b => {
        const latexFormula = b.innerHTML
        let html = ""
        try {
		    html = katex.renderToString(latexFormula, {
			    displayMode: true,
			    strict: false,
			    throwOnError: true,
			    output: 'html',
			    fleqn: true,
			    trust: trustFn,
		    });
	    } catch (e) {
            console.error(e);
		    if (e instanceof katex.ParseError) {
			    html = `<div class="error">Error parsing LaTeX</div>`;
		    }
	    };

        b.innerHTML = html

    })

}

document.addEventListener("DOMContentLoaded", function() {
    const initFns = [initToggles, initLatex]
    initFns.forEach(f => {
        try {
            f()
        } catch (e) {
            console.error(`error executing init function "${f.name}":`, e)
        }
    })
});
