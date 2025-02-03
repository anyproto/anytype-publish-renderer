import $ from 'jquery';
import Prism from 'prismjs';
import mermaid from 'mermaid';
import { instance as viz } from '@viz-js/viz';
import UtilCommon from './lib/common';
import UtilPrism from './lib/prism';

import 'katex/dist/katex.min.css';
import 'prismjs/themes/prism.css';
import 'scss/index.scss';

const katex = require('katex');
require('katex/dist/contrib/mhchem');

for (const lang of UtilPrism.components) {
	require(`prismjs/components/prism-${lang}.js`);
};

declare global {
	interface Window {
		svgSrc: any;
		fathom: any;
	}
};

window.svgSrc = {};

function initToggles () {
    $('.block.textToggle').each((i, block) => {
		block = $(block);
		block.off('click').on('click', () => block.toggleClass('isToggled'));
    });
};

function initLatex () {
    const blocks = $('.block.blockEmbed.isLatex > .content');
    const trustFn = context => [ '\\url', '\\href', '\\includegraphics' ].includes(context.command);

    blocks.each((i, block) => {
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

function initInlineLatex () {
	const blocks = $('.block.blockText > .content > .flex > .text');
	
	blocks.each((i, block) => {
		block = $(block);
		block.html(UtilCommon.getLatex(block.text()));
	});
};

function initMermaid () {
    mermaid.initialize({ startOnLoad: true });
};

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
    });
};

function initPrism () {
	const blocks = $('code');

	blocks.each((i, block) => {
		block = $(block);

		const lang = block.data('lang');
		const value = block.text();

		block.html(Prism.highlight(value, Prism.languages[lang], lang));
	});
};

function initAnalyticsEvents () {
	$('.fathom').each((item, i) => {
		item = $(item);

		item.off('click').on('click', () => {
			window.fathom.trackEvent(item.data('event'));
		});
	});
};

document.addEventListener("DOMContentLoaded", function() {
    const initFns = [ 
		initAnalyticsEvents, 
		initToggles, 
		initLatex, 
		initMermaid, 
		initGraphviz, 
		initPrism, 
		initInlineLatex,
	];

    initFns.forEach(f => {
        setTimeout(_ => {
            try {
                f();
            } catch (e) {
                console.error(`error executing init function "${f.name}":`, e);
            };
        });
    });
});