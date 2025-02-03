const path = require('path');
const process = require('process');
const rspack = require('@rspack/core');
const { RsdoctorRspackPlugin } = require('@rsdoctor/rspack-plugin');

module.exports = (env, argv) => {
	return {
		mode: 'development',
		devtool: 'source-map',

		optimization: {
			minimize: false,
			removeAvailableModules: true,
			removeEmptyChunks: true,
			splitChunks: false,
		},
		
		entry: {
			app: { 
				import: './src/js/entry.js', 
				filename: 'js/main.js',
			},
		},

		 entry: './src/js/entry.js',

		output: {
			path: path.resolve(__dirname, 'static', 'js'),
			filename: 'main.js',
		},

		resolve: {
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
		].filter(Boolean),
	};
};