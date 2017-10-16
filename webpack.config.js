const ExtractTextPlugin = require( 'extract-text-webpack-plugin' ),
      HtmlWebpackPlugin = require( 'html-webpack-plugin' )
      webpack           = require( 'webpack' ),
      plugins           = [

      // omit import xxx
      new webpack.ProvidePlugin({
        React    : 'react',
        ReactDOM : 'react-dom',
        Notify   : 'notify',
        jQuery   : 'jquery',
      }),

      // chunk files
      new webpack.optimize.CommonsChunkPlugin({
        names     : [ 'vendors' ],
        minChunks : Infinity
      }),

      // defined environment variable
      new webpack.DefinePlugin({
        'process.env.NODE_ENV': JSON.stringify( 'production' ) // or development
      }),

      // extract css files
      new ExtractTextPlugin( '[name].css' ),

      // minify html files
      new HtmlWebpackPlugin({
        filename: 'index.html',
        template: 'src/index.html',
        inject: false,
        minify: {
          collapseWhitespace: true,
        },
      }),

    ],

    // conditions environment
    isProduction = function () {
      return process.env.NODE_ENV === 'production';
    },

    // only when environment variable is 'development' call
    develop = ( function () {
      const OpenBrowserPlugin  = require('open-browser-webpack-plugin');
      if ( !isProduction() ) {
        plugins.push(
          new webpack.HotModuleReplacementPlugin(),
          new OpenBrowserPlugin({ url: 'http://localhost:8080' })
        );
      }
    })(),

    // only when environment variable is 'production' call
    deploy = ( function () {
      const CopyWebpackPlugin  = require( 'copy-webpack-plugin'  ),
            CleanWebpackPlugin = require( 'clean-webpack-plugin' );

      // environment verify
      if ( isProduction() ) {

        // delete publish folder
        plugins.push(
          new CleanWebpackPlugin([ 'publish' ], {
            verbose: true,
            dry    : false,
          })
        );

        // copy files
        plugins.push(
          new CopyWebpackPlugin([
            { context: 'src/assets/images/',  from : '*' , to : './assets/images'  },
            { context: 'src/assets/favicon/', from : '*' , to : './assets/favicon' },
          ])
        );

        // call uglifyjs plugin
        plugins.push(
          new webpack.optimize.UglifyJsPlugin({
            compress: {
              sequences: true,
              dead_code: true,
              conditionals: true,
              booleans: true,
              unused: true,
              if_return: true,
              join_vars: true,
              drop_console: true
            },
            mangle: {
              except: [ '$', 'exports', 'require' ]
            },
            output: {
              comments: false
            }
          })
        );

      }
    })(),

    bundle = ( function () {
      const files = [
        './src/index.jsx'
      ];
      if ( !isProduction() ) {
        files.push(
          'webpack/hot/dev-server',
          'webpack-dev-server/client?http://localhost:8080'
        );
      }
      return files;
    }),

    // webpack config
    config = {
      entry: {

        vendors : [

          // react
          './node_modules/react/dist/react.min.js',
          './node_modules/react-dom/dist/react-dom.min.js',

          // vendors
          'jquery',
          'pangu',
          'velocity',

          'wavess',
          'notify',

          // component
          'textfield',
          'button',
          'selectfield',
          /*
          'fab',
          'switch',
          'tabs',
          'sidebar',
          'list',
          'dialog',
          */
          'tooltip',
          'waves'
        ],

        bundle: bundle(),

      },

      output: {
        path     :  isProduction() ? './publish/' : './',
        filename : '[name].js'
      },

      devServer: {
        contentBase: './src',
        port: 8080,
        historyApiFallback: true,
        hot: true,
        inline: true,
        progress: true,
      },

      plugins: plugins,

      module: {
        loaders: [
          {
              test: /\.js[x]?$/,
              exclude: /node_modules/,
              loader: 'babel',
              query: {
                presets: [ 'es2015', 'stage-0', 'react' ]
              }
          },

          // css in js
          //{ test: /\.css$/,         loader: 'style!css!postcss' },

          // extract css files
          { test: /\.css$/,           loader: ExtractTextPlugin.extract( 'style', 'css!postcss' ) },

          // image in js
          { test: /\.(png|jpg|gif)$/, loader: 'url?limit=12288'   },

          // expose $
          {
            test  : require.resolve( './src/vender/jquery-2.1.1.min.js' ),
            loader: 'expose?jQuery!expose?$'
          },

        ]
      },

      postcss: function () {
        return [
          require( 'import-postcss'  )(),
          require( 'postcss-cssnext' )()
        ]
      },

      resolve: {
        alias : {
          jquery     : __dirname + '/src/vender/jquery-2.1.1.min.js',
          pangu      : __dirname + '/src/vender/pangu.min.js',
          velocity   : __dirname + '/src/vender/velocity.min.js',

          wavess     : __dirname + '/src/vender/waves/waves.js',
          notify     : __dirname + '/src/vender/notify/notify.js',

          textfield  : __dirname + '/src/vender/mduikit/textfield.jsx',
          fab        : __dirname + '/src/vender/mduikit/fab.jsx',
          button     : __dirname + '/src/vender/mduikit/button.jsx',
          selectfield: __dirname + '/src/vender/mduikit/selectfield.jsx',
          switch     : __dirname + '/src/vender/mduikit/switch.jsx',
          tabs       : __dirname + '/src/vender/mduikit/tabs.jsx',
          sidebar    : __dirname + '/src/vender/mduikit/sidebar.jsx',
          list       : __dirname + '/src/vender/mduikit/list.jsx',
          dialog     : __dirname + '/src/vender/mduikit/dialog.jsx',
          tooltip    : __dirname + '/src/vender/mduikit/tooltip.jsx',
          waves      : __dirname + '/src/vender/mduikit/waves.js',

          index      : __dirname + '/src/index.jsx',
          entry      : __dirname + '/src/module/entry.jsx',
          search     : __dirname + '/src/module/search.jsx',
          filter     : __dirname + '/src/module/filter.jsx',
          version    : __dirname + '/src/module/version.js',
          controlbar : __dirname + '/src/module/controlbar.jsx',

        }
      }

};

module.exports = config;
