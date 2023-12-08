class Netfetch < Formula
    desc "CLI tool to scan for network policies in Kubernetes clusters/namespaces and provide a score based on the amount of untargeted workloads"
    homepage "https://github.com/deggja/netfetch"
    url "https://github.com/deggja/netfetch/releases/download/0.0.53/netfetch_0.0.53_darwin_amd64.tar.gz"
    sha256 "4bc8092b5df77ce87c6dafa3d46e84779490985dbc2317072b00dd5bafeee61d"
  
    def install
      bin.install "netfetch"
    end
  
    test do
      system "#{bin}/netfetch", "--version"
    end
  end