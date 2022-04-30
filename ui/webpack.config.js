const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const SpriteSmithPlugin = require('webpack-spritesmith')
const TerserPlugin = require('terser-webpack-plugin')

module.exports = {
  module: {
    rules: [
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader']
      },
      {
        test: /\.png$/,
        type: 'asset/resource'
      },
      {
        test: /\.svg$/,
        use: 'svg-url-loader'
      },
      {
        test: /\.js$/,
        include: path.resolve(__dirname, './src'),
        use: 'babel-loader'
      }
    ]
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, 'public/index.html')
    }),
    new SpriteSmithPlugin({
      src: {
        cwd: path.resolve(__dirname, 'src/images/flags'),
        glob: '*.png'
      },
      target: {
        image: path.resolve(__dirname, 'src/images/flags.png'),
        css: path.resolve(__dirname, 'src/images/flags.css')
      }
    })
  ],
  optimization: {
    minimizer: [
      new TerserPlugin({
        extractComments: false
      })
    ]
  },
  output: {
    clean: true,
    filename: '[name].[contenthash].js'
  },
  performance: {
    hints: false,
  },
  devServer: {
    historyApiFallback: true,
    proxy: [
      {
        context: ['/api', '/auth'],
        target: 'http://localhost',
        changeOrigin: true,
        ws: true
      }
    ]
  }
}
