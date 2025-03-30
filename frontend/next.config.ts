import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  output: "standalone",
  //trailingSlash: true,
  //Skip trailing slash redirect to prevent issues with asset paths
  //skipTrailingSlashRedirect: true,
  //Set distDir to 'dist' to output the exported files to the 'dist' folder
  // distDir: "dist",
};

export default nextConfig;
