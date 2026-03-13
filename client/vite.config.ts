import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import type { MinifyOptions } from 'terser';

const terserOptions: MinifyOptions = {
  compress: {
    drop_console: true,
    drop_debugger: true,
    pure_funcs: ['console.log'],
  },
  mangle: true,
  format: {
    comments: false,
  },
};

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  build: {
    minify: 'terser',
    target: 'esnext',
    sourcemap: false,
    terserOptions,
  }
});
