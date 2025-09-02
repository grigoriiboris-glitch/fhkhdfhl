const path = require('path')
const isDev = process.env.NODE_ENV === 'development'


module.exports = {
  lintOnSave: false,
  outputDir: './dist',
  productionSourceMap: false,
  transpileDependencies: ['simple-mind-map'],
  configureWebpack: {
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src/')
      }
    },
    output: {
      hashFunction: "sha256"
    }
  },
  devServer: {
    proxy: {
      '/auth': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        secure: false
      },
      '/api': {
        target: 'http://localhost:8000',
        changeOrigin: true,
        secure: false
      }
    }
  },
  
  publicPath: isDev ? '' : '/hyy-vue3-mindmap/',
  //assetsDir: 'static'
}