
(function() {
	const loader = document.getElementById('root-loader');
	const anim = loader.getElementsByClassName('anim')[0];
	const chunks = %CHUNKS%;
	const length = chunks.length;

	let n = 0;

	const hide = () => {
		loader.remove();
	};

	const loaded = () => {
		requestAnimationFrame(() => anim.classList.remove('from'));

		window.addEventListener('message', e => window.onMessage(e.data));

		window.setTimeout(() => {
			anim.classList.add('to');

			window.setTimeout(() => {
				loader.classList.add('hide');
				window.setTimeout(() => hide(), 300);
			}, 450);
		}, 500);
	};

	chunks.forEach((chunk) => {
		const script = document.createElement('script');

		script.src = '/static/js/build/' + chunk;
		script.defer = true;

		script.onload = function() {
			n++;

			if (n === length) {
				loaded();
			};
		};

		script.onerror = function() {
			console.error(`Error loading chunk: ${chunk}`);
		};

		document.head.appendChild(script);
	});
})();