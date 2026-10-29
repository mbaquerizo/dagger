import { render, screen } from '@testing-library/react';
import type { ComponentType } from 'react';
import { describe, expect, it } from 'vitest';

import { Route } from '#/routes/index';

describe('placeholder route', () => {
  it('renders the placeholder text', () => {
    const App = Route.options.component;

    expect(App).toBeTypeOf('function');

    const Placeholder = App as ComponentType;
    render(<Placeholder />);

    expect(screen.getByText('Hello Dagger.')).toBeInTheDocument();
  });
});
