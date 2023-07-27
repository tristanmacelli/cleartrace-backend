module.exports = {
    root: true,
  
    parser: "@typescript-eslint/parser",
    extends: ["plugin:@typescript-eslint/recommended"],

    parserOptions: {
      ecmaVersion: "latest",
      parser: "@typescript-eslint/parser",
      exclude: ["node_modules"],
    },
  
    plugins: ["prettier", "@typescript-eslint"],
  };
  