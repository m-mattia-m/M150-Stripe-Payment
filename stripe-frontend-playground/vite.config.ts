import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tscPlugin from 'vite-plugin-tsc';

// https://vitejs.dev/config/
export default defineConfig({
  base: '/',
  server: {
    port: 5173,
  },
  plugins: [vue(), tscPlugin()],
  logLevel: 'silent',
});
