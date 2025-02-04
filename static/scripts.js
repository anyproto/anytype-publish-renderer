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
        const latexFormula = b.innerText
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

function initMermaid() {
    mermaid.initialize({ startOnLoad: true });
}

function initGraphviz() {
    const gphBlocks = document.querySelectorAll(".isGraphviz");
    gphBlocks.forEach(b => {
        const gphFormula = window.svgSrc[b.id].content
        try {
            const viz = new Viz()
            viz.renderSVGElement(gphFormula).then(svg => {
                parent = b.querySelector(".content")
                parent.appendChild(svg);
            }, err => {
                console.error("viz error:",err)
            });
        } catch (e) {
            console.error("viz error:",e);
        };
    })
}

function initAnalyticsEvents() {
    document.getElementById("madeInAnytypeLink")?.addEventListener("click", (e) => {
        setTimeout(_ => {
            window.fathom?.trackEvent("PublishSiteClick");
        })
    });

    document.getElementById("joinSpaceLink")?.addEventListener("click", (e) => {
        setTimeout(_ => {
            window.fathom?.trackEvent("PublishJoinSpaceClick");
        })
    });

}
/* When the user clicks on the button,
toggle between hiding and showing the dropdown content */
function myFunction() {
  document.getElementById("myDropdown").classList.toggle("show");
}

// Close the dropdown menu if the user clicks outside of it
window.onclick = function(event) {
  if (!event.target.matches('.dropbtn')) {
    var dropdowns = document.getElementsByClassName("dropdown-content");
    var i;
    for (i = 0; i < dropdowns.length; i++) {
      var openDropdown = dropdowns[i];
      if (openDropdown.classList.contains('show')) {
          openDropdown.classList.remove('show');
      }
    }
  }
}
function initDropdown() {
    document.getElementById("menuButton")?.addEventListener("click", (e) => {
        document.getElementById("menuContent").classList.toggle("show");
    });

    window.addEventListener("click", (e) => {
        if (!event.target.matches('.menuButton img')) {
            document.getElementById("menuContent")?.classList.remove('show');
        }
    })
}

window.svgSrc = {}

document.addEventListener("DOMContentLoaded", function() {
    const initFns = [initDropdown, initAnalyticsEvents, initToggles, initLatex, initMermaid, initGraphviz]
    initFns.forEach(f => {
        setTimeout(_ => {
            try {
                f()
            } catch (e) {
                console.error(`error executing init function "${f.name}":`, e)
            }
        })
    })
});
