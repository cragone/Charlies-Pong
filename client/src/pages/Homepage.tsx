import React, { useState } from "react";
import PingPong from "./../components/PingPong";

const Homepage = () => {
  const [playPingPong, setPlayPingPong] = useState(false);

  if (playPingPong) {
    return <PingPong setOpen={setPlayPingPong} />;
  }

  return (
    <div className="min-h-screen bg-base-200 text-base-content">
      <div className="container mx-auto p-6 space-y-8">
        {/* Header */}
        <div className="text-center">
          <h1 className="text-4xl font-bold">Charles Ragone</h1>
          <p className="text-lg text-gray-400">
            Software Engineer • Go & React Dev
          </p>
        </div>
        {/* Resume Card */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">Resume</h2>
            <p>
              I build fast, concurrent systems in Go and modern web apps in
              React. I’m passionate about solving hard problems with clean code.
            </p>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-4">
              <div>
                <h3 className="font-bold">Languages</h3>
                <p>Go, JavaScript/TypeScript, SQL, Python, Bash</p>
              </div>
              <div>
                <h3 className="font-bold">Tools</h3>
                <p>
                  Docker, Kubernetes, Helm, PostgreSQL, GitHub Actions, AWS,
                  Azure, Terraform
                </p>
              </div>
            </div>
            <div className="card-actions justify-end mt-4">
              <a
                href="/resume.pdf"
                download
                className="btn btn-outline btn-primary"
              >
                Download Resume
              </a>
            </div>
          </div>
        </div>
        {/* Games */}
        <div>
          <h2 className="text-2xl font-bold mb-4">Play Ping Pong</h2>
          <div className="flex justify-center">
            <button
              className="btn btn-primary hover:animate-pulse"
              onClick={() => {
                setPlayPingPong(true);
              }}
            >
              Click Me to Play!!
            </button>
          </div>
        </div>
        {/* Footer */}
        <footer className="mt-12 text-center text-sm text-gray-500">
          Built with ❤️ enjoy!
        </footer>
      </div>
    </div>
  );
};

export default Homepage;
