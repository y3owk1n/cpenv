class Myproject < Formula
  desc "A CLI for copy and paste your local .env to right projects faster"
  homepage "https://github.com/y3owk1n/cpenv"
  
  # You'll replace these in your actual implementation
  version "2.0.0"
  
  # For macOS Intel (x86_64)
  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/y3owk1n/cpenv/releases/download/v#{version}/cpenv-darwin-x64"
    sha256 "f8c2d9b7e9c6fda929cb0b9fba0da9a028d1d3f740213d3d79ae757177642087"
  end
  
  # For macOS Apple Silicon (arm64)
  if OS.mac? && Hardware::CPU.arm?
    url "https://github.com/y3owk1n/cpenv/releases/download/v#{version}/cpenv-darwin-arm64"
    sha256 "844dda394d59eb32854c3c677c68868a1105c9fc32388a650807ccc9c51a2f82"
  end

  def install
    # Rename the binary to match your project name
    bin.install "cpenv-darwin-#{Hardware::CPU.arch}" => "cpenv"
  end

  test do
    # Add a simple test to verify installation
    system "#{bin}/cpenv", "--version"
  end
end
