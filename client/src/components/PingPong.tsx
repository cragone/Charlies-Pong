import React, { useState, useEffect } from "react";

const PingPong = ({ setOpen }) => {
  useEffect(() => {
    const handleKeyDown = (event) => {
      if (event.ctrlKey && event.key.toLowerCase() === "c") {
        setOpen(false);
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [setOpen]);

  return (
    <>
      <div className="bg-gray-900 text-success rounded-lg shadow-lg overflow-hidden font-mono w-full max-w-2xl">
        {/* Title bar */}
        <div className="bg-gray-800 flex items-center px-3 py-2">
          <div className="flex space-x-2">
            <span
              className="w-3 h-3 bg-error rounded-full"
              onClick={() => {
                setOpen(false);
              }}
            ></span>
            <span className="w-3 h-3 bg-warning rounded-full"></span>
            <span className="w-3 h-3 bg-success rounded-full"></span>
          </div>
          <span className="ml-4 text-success text-sm">Ping Pong</span>
        </div>

        {/* Terminal content */}
        <div className="p-4 space-y-1">
          <div className="whitespace-pre-wrap">
            <span className="text-success">$ </span>
            <span>Hello World!</span>
          </div>
          {/* Blinking cursor */}
          <div>
            <span className="text-success">$ </span>
            <span className="animate-pulse">█</span>
          </div>
        </div>
      </div>
    </>
  );
};

export default PingPong;
