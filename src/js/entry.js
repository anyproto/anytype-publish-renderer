import $ from 'jquery';
import Prism from 'prismjs';
import mermaid from 'mermaid';
import { instance as viz } from '@viz-js/viz';

import 'katex/dist/katex.min.css';
import 'prismjs/themes/prism.css';
import 'scss/index.scss';

const katex = require('katex');
require('katex/dist/contrib/mhchem');

window.svgSrc = {};

function initToggles () {
    const blocks = $('.textToggle');

    blocks.each(block => {
		block = $(block);

		block.off('click').on('click', () => {
			block.classToggle('isToggled');
		});
    });
};

function initLatex () {
    const blocks = $('.isLatex .content');
    const trustFn = context => [ '\\url', '\\href', '\\includegraphics' ].includes(context.command);

    blocks.each(block => {
		block = $(block);

        let html = '';
        try {
            html = katex.renderToString(block.text(), {
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
            };
        };

        block.html(html);
    });
};

function initMermaid () {
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

function initAnalyticsEvents () {
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
};

document.addEventListener("DOMContentLoaded", function() {
    const initFns = [initAnalyticsEvents, initToggles, initLatex, initMermaid, initGraphviz]
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