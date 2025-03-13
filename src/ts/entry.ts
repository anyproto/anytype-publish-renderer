import $ from 'jquery/dist/jquery.js';
import raf from 'raf';
import Prism from 'prismjs';
import mermaid from 'mermaid';
import { instance as viz } from '@viz-js/viz';
import * as pdfjs from 'pdfjs-dist';
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

pdfjs.GlobalWorkerOptions.workerSrc = require('pdfjs-dist/build/pdf.worker.js');

declare global {
	interface Window {
		svgSrc: any;
		fathom: any;
		CoverParam: any;
		onMessage: any;
	}
};

window.svgSrc = {};

function renderCover () {
	const block = $('.block.blockCover');
	if (!block.length) {
		return;
	};

	const { CoverX: x, CoverY: y, CoverScale: scale } = window.CoverParam || {};
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

function renderToggles () {
    $('.block.textToggle').each((i, block) => {
		block = $(block);
		block.find('> .content .marker.toggle').off('click').on('click', function () {
			block.toggleClass('isToggled');
		});
    });
};

function renderLatex () {
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

function renderInlineLatex () {
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

function renderMermaid () {
	mermaid.initialize({ 
		securityLevel: 'loose',
		theme: 'base', 
		themeVariables: {
			fontFamily: 'Helvetica, Arial, sans-serif',
			fontSize: '14px'
		},
		startOnLoad: true,
	});

	mermaid.run({ querySelector: `.mermaidChart` });
};

function renderGraphviz () {
	const blocks = $(`.block.blockEmbed.isGraphviz > .content`);

	blocks.each((i, block) => {
		block = $(block);

		viz().then(viz => {
			const text = block.text();
			block.html(viz.renderSVGElement(text));
		});
	});
};

function renderPrism () {
	const blocks = $('.block.blockText.textCode > .content > .flex > .text');

	blocks.each((i, block) => {
		block = $(block);

		const lang = block.data('lang') || 'plain';
		const value = block.text();

		block.html(Prism.highlight(value, Prism.languages[lang], lang));
	});
};

function renderAnalyticsEvents () {
	$('.fathom').each((i, item) => {
		item = $(item);

		item.off('click').on('click', () => {
			window.fathom.trackEvent(item.data('event'));
		});
	});
};

function renderPdf () {
	const blocks = $('.block.blockMedia.isPdf > .content > .wrap');

	let page = 1;

	blocks.each((i, block) => {
		block = $(block);

		const { id, src } = block.data();
		const loadingTask = pdfjs.getDocument(src);

		loadingTask.promise.then(pdf => {
			const switchPage = (page) => {
				const pager = block.find('.pager');
				const number = pager.find('.number');
				const arrowLeft = pager.find('.arrow.left');
				const arrowRight = pager.find('.arrow.right');
				const arrowEndLeft = pager.find('.arrow.end.left');
				const arrowEndRight = pager.find('.arrow.end.right');
				const canLeft = page > 1;
				const canRight = page < pdf.numPages;

				if (pdf.numPages == 1) {
					pager.hide();
				};

				number.text(`${page} / ${pdf.numPages}`);

				arrowLeft.toggleClass('disabled', !canLeft);
				arrowEndLeft.toggleClass('disabled', !canLeft);
				arrowRight.toggleClass('disabled', !canRight);
				arrowEndRight.toggleClass('disabled', !canRight);

				if (canLeft) {
					arrowLeft.off('click').on('click', () => switchPage(page - 1));
					arrowEndLeft.off('click').on('click', () => switchPage(1));
				};

				if (canRight) {
					arrowRight.off('click').on('click', () => switchPage(page + 1));
					arrowEndRight.off('click').on('click', () => switchPage(pdf.numPages));
				};

				pdf.getPage(page).then((page) => {
					const scale = 1.5;
					const viewport = page.getViewport({ scale });
					const canvas = block.find(`#pdfCanvas-${id}`).get(0);
					const context = canvas.getContext('2d');

					canvas.width = viewport.width;
					canvas.height = viewport.height;

					const renderContext = {
						canvasContext: context,
						viewport: viewport,
					};

					page.render(renderContext);
				});
			};

			switchPage(page);
		}).catch((error) => {
			console.error("Error loading PDF:", error);
		});
	});
};

function renderMenu () {
	const win = $(window);
	const btn = $('.menuButton');

	btn.off('click').on('click', function () {
		const el = $(this);
		const st = win.scrollTop();
		const { menuId, horizontal } = el.data();
		const wrap = $(`#menu-${menuId}`);
		const dimmer = wrap.find('.dimmer');
		const menu = wrap.find('.menu');
		const css: any = {};

		wrap.show();

		const { left, top } = el.offset();
		const ew = el.width();
		const eh = el.height();
		const iw = menu.width();

		css.top = top - st + eh + 4;

		switch (horizontal) {
			case 'left': {
				css.left = left;
				break;
			};

			case 'right': {
				css.left = left + ew - iw;
				break;
			};

			case 'center': {
				css.left = left + (ew - iw) / 2;
				break;
			};
		};

		menu.css(css);
		raf(() => {
			el.addClass('hover');
			wrap.addClass('show');
		});

		const hide = () => {
			el.removeClass('hover');
			wrap.removeClass('show');

			setTimeout(() => wrap.hide(), 200);
			win.off('keydown');
		};

		win.off('keydown').on('keydown', e => {
			const k = e.key.toLowerCase();

			if (k == 'escape') {
				hide();
			};
		});

		dimmer.off('click').on('click', () => hide());
	});
};

function renderTable () {
	const blocks = $('.block.blockTable');

	blocks.each((i, block) => {
		block = $(block);

		const wrapper = $('.blocks');
		const scrollWrap = block.find('.scrollWrap');
		const table = block.find('.table');	
		const parent = block.parents('.block');
		const wrapperWidth = wrapper.width();

		let width = table.width();
		let maxWidth = 0;

		if (!parent.length) {
			maxWidth = $(window).width() -  92;
			width = Math.max(wrapperWidth, Math.min(maxWidth, width));
			block.css({
				width: (width >= wrapperWidth) ? width : 'auto',
				marginLeft: (width >= wrapperWidth) ? Math.min(0, (wrapperWidth - width) / 2) : '',
			});
		} else {
			maxWidth = parent.width();
		};

		scrollWrap.toggleClass('withScroll', width >= maxWidth);
	});
};

function renderRow () {
	const wrapper = $('.blocks');
	const blocks = $('.block.blockLayout.layoutRow');
	const mw = wrapper.width() / 2;

	blocks.each((i, block) => {
		block = $(block);

		const children = block.find('> .children > .block');
		const length = children.length;

		children.each((i, child) => {
			child = $(child);

			const width = (Number(child.data('width')) || 1 / length) * 100;
			const innerBlocks = child.find('> .children > .block');
			
			child.css({ width: `${width}%` });

			raf(() => {
				innerBlocks.toggleClass('isVertical', child.width() < mw);
			});
		});
	});
};

function renderColumn () {
	const blocks = $('.block.blockLayout.layoutColumn');

	blocks.each((i, block) => {
		block = $(block);

		const children = block.find('> .children > .block');

		children.first().addClass('first');
		children.last().addClass('last');
	});
};

function renderPrimitivesSvgs () {
    document.querySelectorAll(".svg-container").forEach(container => {
        const svgUrl = container.getAttribute("data-src");
		console.log("svgUrl:", svgUrl);

        const color = container.getAttribute("data-color");

        if (!svgUrl || !color) return;

        fetch(svgUrl)
            .then(response => response.text())
            .then(svgText => {
                const updatedSvg = svgText.replace(/_COLOR_VAR_/g, color);
                const parser = new DOMParser();
                const svgDoc = parser.parseFromString(updatedSvg, "image/svg+xml").documentElement;

				// some svgs in assets have width and height set to e.g. 512
				// redefine to 100% to respect the parrent
				svgDoc.attributes.width.value = "100%"
				svgDoc.attributes.height.value = "100%"


                container.innerHTML = "";
                container.appendChild(svgDoc);
            })
            .catch(error => console.error("Error loading SVG:", error));
    });
}

window.onMessage = (data) => {
	const { type, height, blockId, url } = data;

	switch (type) {
		case 'resize': {
			$(`#receiver${blockId}`).css({ height: Math.max(80, height) });
			break;
		};

		case 'openUrl': {
			window.open(url, '_blank');
			break;
		};
	};
};

$(document).ready(() => {
	const win = $(window);

	win.off('resize').on('resize', () => { 
		renderCover();
		renderTable();
		renderRow();
	});

    const renderFns = [ 
		renderTable,
		renderRow,
		renderColumn,
		renderCover,
		renderAnalyticsEvents, 
		renderToggles, 
		renderLatex, 
		renderMermaid, 
		renderGraphviz, 
		renderPrism, 
		renderInlineLatex,
		renderPdf,
		renderMenu,
		renderPrimitivesSvgs,
	];

	renderFns.forEach(f => {
		setTimeout(_ => {
			try {
				f();
			} catch (e) {
				console.error(`error executing render function "${f.name}":`, e, f);
			};
		});
	});
});
