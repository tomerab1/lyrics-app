import React from "react";

export function Card({ title, children }: { title: React.ReactNode; children: React.ReactNode }) {
  return (
    <section className="bg-white/90 backdrop-blur rounded-2xl shadow-lg ring-1 ring-black/5 p-5 md:p-7 transition-shadow hover:shadow-xl">
      <div className="mb-4">
        {typeof title === 'string' ? (
          <h2 className="text-lg md:text-xl font-semibold text-center text-gray-900">{title}</h2>
        ) : (
          title
        )}
      </div>
      {children}
    </section>
  );
}


