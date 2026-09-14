import { defineConfig } from 'astro/config';
import tailwindcss from '@tailwindcss/vite';

// Static marketing front door: zero JS by default, served as dist/
// on its own Railway service (root dir /web).
export default defineConfig({
  output: 'static',
  vite: {
    plugins: [tailwindcss()],
  },
});
