class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters and provide a score"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.28/netfetch_0.0.28_darwin_amd64.tar.gz"
    sha256 "df9877174d0fcdc733a1c7af9fe4cc6ae99b31b2ae0a6fe2eee87d8c5d1bffc3"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end