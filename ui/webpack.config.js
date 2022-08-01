const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const SpriteSmithPlugin = require('webpack-spritesmith')
const TerserPlugin = require('terser-webpack-plugin')

function generateSpriteStylesheet(data) {
  const s = data.sprites[0]
  const generic = `
.icon {
  background-image: url(${s.image});
  background-size: ${s.total_width / 2}px ${s.total_height / 2}px;
  height: ${s.height / 2}px;
  width: ${s.width / 2}px;
}
  `
  const specific = data.sprites.map((s) => `
.icon-${s.name} {
  background-position: ${s.offset_x / 2}px ${s.offset_y / 2}px;
}
  `).join('')
  return generic + specific
}

module.exports = {
  module: {
    rules: [
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader', 'postcss-loader']
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
        css: [
          [path.resolve(__dirname, 'src/images/flags.css'), {
            format: 'default'
          }]
        ]
      },
      customTemplates: {
        'default': generateSpriteStylesheet
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
    filename: '[name].[contenthash].js',
    publicPath: '/'
  },
  performance: {
    hints: false,
  },
  devtool: 'nosources-source-map',
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
