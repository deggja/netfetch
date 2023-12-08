class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.55/netfetch_0.0.55_darwin_amd64.tar.gz"
    sha256 "5ca27a17e1c5d89943c9a410b6b7baf460dd3f8ce792be70923ceb2231b4cd66"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end