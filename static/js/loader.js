
(function() {
	const loader = document.getElementById('root-loader');
	const anim = loader.getElementsByClassName('anim')[0];
	const chunks = ["main~13.fa6a8969adc85dc0.js","6153.b9f0f9a49e16aa8c.js","main~0.9aa0c72ea0729ca6.js","main~1.ea38bbbc4efcd3fe.js","main~2.36ecb7ae21e4a113.js","main~3.47694bbb3207b746.js","main~4.a9e1fe0f067e2cd8.js","main~5.d11a402a81d49b63.js","main~6.e3347157b6051763.js","main~7.eeb10caaaa13d742.js","main~8.8d401f0ad24e4555.js","main~9.b62d5307f0332171.js","main~10.a01447b856e083de.js","main~11.ab7944c51213fe7d.js","main~12.20684202d4202ae4.js","9658.650b2084b4423a59.js","1746.67589837305d3bb3.js","6667.120fc73e9c94ef6b.js","4338.1bd85e0f7c2e4104.js","6145.38944b42d76483ee.js","4137.ced8cb52614c21a6.js","6211.f3bf22f5f99a32b0.js","4938.2d6fb1d0a8f49e96.js","7961.21a4aa76607e0e44.js","9574.e291b77be957b3bb.js","4709.c1aa7321c21cc0ad.js","829.78cd56bde3927f2f.js","1961.d6c627fdfefe1dfd.js","8260.0352716e29cb957d.js","7044.0b917a7fb95bb9cd.js","4524.98197999e698c410.js","6622.a9fba4249a76b927.js","6972.13cd2308d6f16648.js","2195.710e8710a5406fbe.js","8734.1eeb80ecc9ec7db2.js","7192.e5002661fdd84110.js","565.6df5700dd4636f3a.js","7994.e34c88e96b9fe8c6.js","2101.1cb9f034ae933735.js","1704.99aaa6acc507c50c.js","5146.e9e6f746e1fdf8b0.js","8720.4ab83474fa957131.js","306.09b702bd1e955fdf.js","6299.dbe11fc88532ecb8.js","1282.afed4f72307e8873.js","287.d22eb566b99ecff7.js","5455.eff7ff3bd55d8ff8.js","562.f480cb1f0be86030.js","6226.6d595f9183dc948f.js","106.3b8bb7db728e5dfb.js","3917.67643c494e1cb7f1.js","362.2422b211b03d4234.js","3337.7eb2f3affa022e80.js","2646.4b1ac2eb85673e2a.js","9351.528d516db2448a5d.js","1194.61d20fdff17718e4.js","5597.a3e74c5e43c332ac.js","1904.d19de85e1a31e365.js","5459.2a12aec72dffbea5.js","7977.924143bd4764274e.js","5905.ec922cac3a4b19ad.js","3315.4c45894ece4f2792.js","2731.c20b6f96e4cc0c3f.js","9757.7e161e5afbd02f48.js","5439.440ca00a9a203a35.js","7558.a51457b93039b661.js","4600.d93663e8fcdf4b04.js","488.d5b95f1e290a5f6d.js","4257.c4eef31e015c020a.js","1477.f58c85d510b845fd.js","8810.ee1f5844121e7b4e.js","6518.f259bab75891ef56.js","7199.409e85887e9ea044.js","2623.0fdaec0b86e3c9c8.js","379.629bf5372fbbfe8c.js","2693.5add90d47171dba2.js","9284.f6c35e2209aed925.js","4017.10dca8071823fd7a.js","5552.d3ad932eba1f9878.js","8680.05cd3be3fc2613a8.js","3878.437df47a869076cf.js","6588.b75794b0b3d06bf4.js","7042.88542415bed12944.js","6046.6084aa2293d50408.js","3720.6135a348b2017d6a.js","966.2cd5299c9da32315.js","5294.e5bb9118a7660bb3.js"];
	const length = chunks.length;

	console.log("Chunks to be loaded:", chunks);

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

			console.log(`Chunk ${chunk} loaded successfully.`);
		};

		script.onerror = function() {
			console.error(`Error loading chunk: ${chunk}`);
		};

		document.head.appendChild(script);
	});
})();