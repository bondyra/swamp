const { override, addWebpackResolve } = require("customize-cra");

module.exports = override(
  addWebpackResolve({
    fallback: {
      crypto: false,
      path: false,
      fs: false,
    },
  })
);
