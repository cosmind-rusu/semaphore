module.exports = {
  configureWebpack: {
    devServer: {
      historyApiFallback: true,
      proxy: {
        '^/api': {
          target: 'http://localhost:3330',
        },
        '^/billing': {
          target: 'http://localhost:2500',
          pathRewrite(path) { return path.replace('/billing/', '/'); },
        },
        '^/helpdesk': {
          target: 'http://localhost:3000',
          pathRewrite(path) { return path.replace('/helpdesk/', '/'); },
        },
      },
    },
  },
  chainWebpack: (config) => {
    config.plugin('html')
      .tap((args) => {
        // eslint-disable-next-line no-param-reassign
        args[0].minify = false;
        return args;
      });
  },
  transpileDependencies: [
    'vuetify',
  ],
  publicPath: './',
  outputDir: '../api/public',
};
