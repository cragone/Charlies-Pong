import React from "react";

const Homepage = () => {
  const projects = [
    {
      title: "CharliesPong",
      description:
        "Terminal-based Ping Pong game in Go using goroutines and real-time input handling.",
      link: "https://github.com/cragone/charliespong",
    },
    // Add more projects here...
  ];

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
                <p>Go, JavaScript/TypeScript, SQL, Bash</p>
              </div>
              <div>
                <h3 className="font-bold">Tools</h3>
                <p>Docker, Kubernetes, Helm, PostgreSQL, GitHub Actions</p>
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
        {/* Portfolio */}
        <div>
          <h2 className="text-2xl font-bold mb-4">Projects</h2>
          <div className="grid gap-4 md:grid-cols-2">Play Pomg Here</div>
        </div>
        \{/* Footer */}
        <footer className="mt-12 text-center text-sm text-gray-500">
          Built with ❤️ using Vite, React, TailwindCSS & DaisyUI
        </footer>
      </div>
    </div>
  );
};

export default Homepage;
