import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/')({ component: App });

function App() {
  return (
    <main>
      <header>
        <h1 className="text-center text-2xl">Dagger</h1>
      </header>
      Hello Dagger.
    </main>
  );
}
