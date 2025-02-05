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
		CoverParam: any;
	}
};

window.svgSrc = {};

function initCover () {
	const { x, y, scale } = window.CoverParam;
	const block = $('.block.blockCover');
	const cover = block.find('#cover');
	const bw = block.width();
	const bh = block.height();

	cover.css({ height: 'auto', width: `${(scale + 1) * 100}%` });

	const cw = cover.width();
	const ch = cover.height();
	const mx = cw - bw;
	const my = ch - bh;

	let newX = x * cw;
	let newY = y * ch;

	newX = Math.max(-mx, Math.min(0, newX));
	newY = Math.max(-my, Math.min(0, newY));

	const px = (newX / cw) * 100;
	const py = (newY / ch) * 100;
	const css: any = { transform: `translate3d(${px}%,${py}%,0px)` };

	if (ch < bh) {
		css.transform = 'translate3d(0px,0px,0px)';
		css.height = bh;
		css.width = 'auto';
	};

	cover.css(css);
};

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
	const blocks = $(`
		.block.blockText:not(.textCode) > .content > .flex > .text,
		.block.blockTableOfContents > .content .item a,
		.block.blockLink > .content .name,
		.block.blockLink > .content .description
	`);
	
	blocks.each((i, block) => {
		block = $(block);
		block.html(UtilCommon.getLatex(block.html()));
	});
};

function initMermaid () {
    mermaid.initialize({ 
		 securityLevel: 'loose',
		theme: 'base', 
		themeVariables: {
			fontFamily: 'Helvetica, Arial, sans-serif',
			fontSize: '14px'
		},
		startOnLoad: true,
	});
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
	const blocks = $('.block.blockText.textCode > .content > .flex > .text');

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
		initCover,
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