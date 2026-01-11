
import React from 'react';
import { render } from '@testing-library/react';
import ScenarioLiveRoom from './ScenarioLiveRoom';
import axios from 'axios';
import { MemoryRouter } from 'react-router-dom';
import { ErrorBoundary } from 'react-error-boundary';

jest.mock('axios');
jest.mock('react-i18next', () => ({
  useTranslation: () => ({ t: key => key }),
}));
jest.mock('../components/LanguageSwitch', () => ({
  useSrsLanguage: () => 'en',
}));
jest.mock('../utils', () => ({
  Token: {
    loadBearerHeader: () => ({}),
  },
  Locale: {
    current: () => 'en',
  },
}));

describe('ScenarioLiveRoom', () => {
  test('renders list of rooms', async () => {
    const mockRooms = [
      { uuid: '1', title: 'Room 1', created_at: '2023-01-01' },
      { uuid: '2', title: 'Room 2', created_at: '2023-01-02' },
    ];

    axios.post.mockResolvedValue({
      data: {
        data: {
          rooms: mockRooms,
        },
      },
    });

    render(
      <MemoryRouter>
        <ErrorBoundary fallback={<div>Error</div>}>
          <ScenarioLiveRoom />
        </ErrorBoundary>
      </MemoryRouter>
    );

    // Check if the component requests room list
    expect(axios.post).toHaveBeenCalledWith('/terraform/v1/live/room/list', {}, expect.any(Object));
  });
});
