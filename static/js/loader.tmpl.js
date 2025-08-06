(function() {
	const chunks = %CHUNKS%;
	const cssFiles = %CSS%;

	cssFiles.forEach(file => {
		const link = document.createElement('link');
		link.rel = 'stylesheet';
		link.href = '/static/js/build/' + file;
		document.head.appendChild(link);
	});

	chunks.forEach(file => {
		const script = document.createElement('script');
		script.src = '/static/js/build/' + file;
		script.defer = true;
		script.onerror = function() {
			console.error(`Error loading chunk: ${file}`);
		};
		document.head.appendChild(script);
	});
})();
