const path = require('path');
const process = require('process');
const rspack = require('@rspack/core');
const { RsdoctorRspackPlugin } = require('@rsdoctor/rspack-plugin');

module.exports = (env, argv) => {
	const port = process.env.SERVER_PORT;

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

		watchOptions: {
			ignored: /node_modules/,
			poll: false,
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