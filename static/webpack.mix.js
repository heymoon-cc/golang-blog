let mix = require('webpack-mix');

// noinspection JSVoidFunctionReturnValueUsed
mix.js('src/js/app.js', 'dist').js('src/js/admin.js', 'dist')
  .sass('src/css/app.scss', 'dist');
