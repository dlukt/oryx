import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import ScenarioDubbing from './ScenarioDubbing';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import axios from 'axios';
import '@testing-library/jest-dom';

// Mock react-i18next
jest.mock('react-i18next', () => ({
  useTranslation: () => ({ t: (key) => key }),
  I18nextProvider: ({ children }) => <div>{children}</div>,
  initReactI18next: { type: '3rdParty', init: jest.fn() },
}));

// Mock react-error-boundary
jest.mock('react-error-boundary', () => ({
  useErrorBoundary: () => ({ showBoundary: jest.fn() }),
  ErrorBoundary: ({ children }) => <div>{children}</div>,
}));

// Mock axios
jest.mock('axios');

// Mock utils
jest.mock('../utils', () => ({
  Token: {
    load: jest.fn(() => ({ token: 'mock-token' })),
    loadBearer: jest.fn(() => ({ token: 'mock-token' })),
    loadBearerHeader: jest.fn(() => ({ Authorization: 'Bearer mock-token' })),
  },
  Locale: {
    current: jest.fn(() => 'en'),
  }
}));

// Mock components
jest.mock('../components/LanguageSwitch', () => ({
  useSrsLanguage: () => 'en',
}));

jest.mock('../components/SrsErrorBoundary', () => ({
  SrsErrorBoundary: ({ children }) => <div>{children}</div>,
}));

jest.mock('../components/OpenAISettings', () => ({
  OpenAISecretSettings: () => <div>OpenAI Settings</div>,
}));

jest.mock('../components/VideoSourceSelector', () => () => <div>Video Source Selector</div>);

// Mock react-bootstrap-icons
jest.mock('react-bootstrap-icons', () => ({
  Soundwave: () => <span>Soundwave</span>,
  Clipboard: () => <span>Clipboard</span>,
}));

describe('ScenarioDubbing Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  test('renders dubbing editor and polls task', async () => {
    const dubbingId = 'project-uuid-123';
    const taskUuid = 'task-uuid-456';

    const mockProject = {
      uuid: dubbingId,
      title: 'My Dubbing Project',
      filepath: 'video.mp4',
      created_at: '2023-01-01',
      format: { duration: 100, bit_rate: 1000000 },
      task: { uuid: taskUuid, status: 'processing' }, // Start with a task
      video: { codec_type: 'video', codec_name: 'h264', width: 1920, height: 1080, profile: 'high', level: '4.0' },
      audio: { codec_type: 'audio', codec_name: 'aac', profile: 'lc', sample_rate: 44100, channels: 2 },
      asr: { aiProvider: 'openai' } // Ensure editor is rendered
    };

    const mockTaskProcessing = {
      uuid: taskUuid,
      status: 'processing',
      asr_response: {
        groups: [{ uuid: 'g1', segments: [{uuid: 's1', start: 0, end: 1, text: 'Hello'}] }]
      }
    };

    axios.post.mockImplementation((url) => {
      if (url === '/terraform/v1/dubbing/query') {
        return Promise.resolve({ data: { data: mockProject } });
      }
      if (url === '/terraform/v1/dubbing/task-start') {
        return Promise.resolve({ data: { data: mockTaskProcessing } });
      }
      if (url === '/terraform/v1/dubbing/task-query') {
        return Promise.resolve({ data: { data: mockTaskProcessing } });
      }
      if (url === '/terraform/v1/mgmt/openai/query') {
        return Promise.resolve({ data: { data: {} } });
      }
      return Promise.resolve({ data: { data: {} } });
    });

    render(
      <MemoryRouter initialEntries={[`/dubbing?dubbingId=${dubbingId}`]}>
          <Routes>
              <Route path="/dubbing" element={<ScenarioDubbing />} />
          </Routes>
      </MemoryRouter>
    );

    // Verify editor is rendered
    await waitFor(() => {
        expect(screen.getByText(/dubb.studio.start/)).toBeInTheDocument();
    });

    // Verify polling behavior calls task-query
    await waitFor(() => {
        expect(axios.post).toHaveBeenCalledWith('/terraform/v1/dubbing/task-query', expect.objectContaining({
            uuid: dubbingId,
            task: taskUuid
        }), expect.any(Object));
    });
  });
});
