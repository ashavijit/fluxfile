# Flux - Modern Task Runner & Build Automation
# Homepage: https://github.com/ashavijit/fluxfile
# Documentation: https://github.com/ashavijit/fluxfile#readme

class Flux < Formula
  desc "Modern task runner and build automation tool with clean, minimal syntax"
  homepage "https://github.com/ashavijit/fluxfile"
  version "VERSION_PLACEHOLDER"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/ashavijit/fluxfile/releases/download/VERSION_PLACEHOLDER/flux-VERSION_PLACEHOLDER-darwin-amd64.tar.gz"
      sha256 "SHA256_DARWIN_AMD64_PLACEHOLDER"
    end
    on_arm do
      url "https://github.com/ashavijit/fluxfile/releases/download/VERSION_PLACEHOLDER/flux-VERSION_PLACEHOLDER-darwin-arm64.tar.gz"
      sha256 "SHA256_DARWIN_ARM64_PLACEHOLDER"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/ashavijit/fluxfile/releases/download/VERSION_PLACEHOLDER/flux-VERSION_PLACEHOLDER-linux-amd64.tar.gz"
      sha256 "SHA256_LINUX_AMD64_PLACEHOLDER"
    end
    on_arm do
      url "https://github.com/ashavijit/fluxfile/releases/download/VERSION_PLACEHOLDER/flux-VERSION_PLACEHOLDER-linux-arm64.tar.gz"
      sha256 "SHA256_LINUX_ARM64_PLACEHOLDER"
    end
  end

  def install
    if OS.mac?
      if Hardware::CPU.intel?
        bin.install "flux-darwin-amd64" => "flux"
      else
        bin.install "flux-darwin-arm64" => "flux"
      end
    elsif OS.linux?
      if Hardware::CPU.intel?
        bin.install "flux-linux-amd64" => "flux"
      else
        bin.install "flux-linux-arm64" => "flux"
      end
    end
  end

  test do
    assert_match "Flux version", shell_output("#{bin}/flux -v")
  end
end
