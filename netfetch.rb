class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters and provide a score"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.35/netfetch_0.0.35_darwin_amd64.tar.gz"
    sha256 "dfbbf121e60641ac5b1bad7abf297510178765e7b2f87049040fe5da076f7143"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end