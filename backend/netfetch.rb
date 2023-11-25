class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.23/netfetch_0.0.23_darwin_amd64.tar.gz"
    sha256 "edda5dd7c8591a4620d928398b3eefb0436206e37496442a8b10359ef0e86ce3"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end