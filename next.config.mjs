/** @type {import('next').NextConfig} */
const nextConfig = {
  // ESLint is intentionally not configured (keeps dependencies lean); type
  // safety is still enforced by the TypeScript build and `npm run typecheck`.
  eslint: { ignoreDuringBuilds: true },

  // In local development the Go API runs as a separate process. Setting
  // API_PROXY_TARGET (e.g. http://localhost:8080) makes Next proxy /api/* to it
  // so the browser stays same-origin and the session cookie works. In production
  // on Vercel, /api/* is served by the Go function, so no proxy is configured.
  async rewrites() {
    const target = process.env.API_PROXY_TARGET;
    if (!target) return [];
    return [{ source: "/api/:path*", destination: `${target}/api/:path*` }];
  },
};

export default nextConfig;
