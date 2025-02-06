const path = require('path');
const fs = require('fs');
const process = require('process');
const rspack = require('@rspack/core');
const { RsdoctorRspackPlugin } = require('@rsdoctor/rspack-plugin');

module.exports = (env, argv) => {
	const prod = argv.mode === 'production';

	return {
		mode: 'development',
		devtool: 'source-map',

		optimization: {
			minimize: true,
			removeAvailableModules: true,
			removeEmptyChunks: true,
			splitChunks: {
				chunks: 'all',
				minSize: 20000,
				maxSize: 200000,
				maxAsyncRequests: 30,
				maxInitialRequests: 30,
				cacheGroups: {
					defaultVendors: {
						test: /[\\/]node_modules[\\/]/,
						priority: -10,
						reuseExistingChunk: true,
					},
					default: {
						minChunks: 2,
						priority: -20,
						reuseExistingChunk: true,
					},
				},
			},
		},
		
		entry: './src/ts/entry.ts',

		output: {
			path: path.resolve(__dirname, 'static', 'js', 'build'),
			filename: '[name].[contenthash].js',
			chunkFilename: '[name].[contenthash].chunk.js',
			clean: true,
		},

		resolve: {
			extensions: [ '.ts', '.tsx', '.js', '.jsx' ],
			alias: {
				dist: path.resolve(__dirname, 'dist'),
			},
			modules: [
				path.resolve('./src/'),
				path.resolve('./dist/'),
				path.resolve('./node_modules')
			]
		},

		module: {
			rules: [
				{
					test: /\.(j|t)s$/,
					exclude: [/[\\/]node_modules[\\/]/],
					loader: 'builtin:swc-loader',
					options: {
						jsc: {
							parser: {
								syntax: 'typescript',
							},
							transform: {
								react: {
									runtime: 'automatic',
									development: !prod,
									refresh: !prod,
								},
							},
						},
						env: {
							targets: 'Chrome >= 48',
						},
					},
				},
				{
					test: /\.(j|t)sx$/,
					loader: 'builtin:swc-loader',
					exclude: [/[\\/]node_modules[\\/]/],
					options: {
						jsc: {
							parser: {
								syntax: 'typescript',
								tsx: true,
							},
							transform: {
								react: {
									runtime: 'automatic',
									development: !prod,
									refresh: !prod,
								},
							},
						},
						env: {
							targets: 'Chrome >= 48', // browser compatibility
						},
					},
				},
				{
					test: /\.(eot|ttf|otf|woff|woff2)$/,
					type: 'asset/inline'
				},
				{
					test: /\.(jpe?g|png|gif|svg)$/,
					type: 'asset/inline'
				},
				{
					test: /\.s?css/,
					use: [
						{ loader: 'style-loader' },
						{ loader: 'css-loader' },
						{ loader: 'sass-loader' }
					]
				}
			]
		},

		plugins: [
			process.env.RSDOCTOR && new RsdoctorRspackPlugin({}),
			new rspack.optimize.LimitChunkCountPlugin({ maxChunks: 1 }),

			// Custom plugin to generate chunk-loader.js after the build
			{
				apply: (compiler) => {
					compiler.hooks.emit.tapAsync('ChunkLoaderPlugin', (compilation, callback) => {
						const chunks = [];

						compilation.chunks.forEach((chunk) => {
							chunk.files.forEach((file) => {
								if (file.endsWith('.js')) {
									chunks.push(file);
								};
							});
						});

						let chunkLoaderContent = fs.readFileSync(path.resolve(__dirname, 'static', 'js', 'loader.tmpl.js'), 'utf8');

						// Replace %CHUNKS% with the actual chunks
						chunkLoaderContent = chunkLoaderContent.replace('%CHUNKS%', JSON.stringify(chunks));

						fs.writeFileSync(path.resolve(__dirname, 'static', 'js', 'loader.js'), chunkLoaderContent);
						callback();
					});
				},
			},
		].filter(Boolean),
	};
};