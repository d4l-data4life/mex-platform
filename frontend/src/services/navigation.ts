import { CONFIG_URL, LANGUAGE_CODES, NAVIGATION_CONFIG } from 'config';
import { get } from 'utils/fetch-client';
import { LinkItem } from 'config/navigation';

interface NavigationResponse {
  header: LinkItem[];
  footer: LinkItem[];
  learningEnvironment: LinkItem[];
  services: LinkItem[];
}

export class NavigationService {
  async fetch(): Promise<NavigationResponse> {
    const [response] = await get<NavigationResponse>({
      url: `${CONFIG_URL}/pages`,
    });

    return response;
  }

  async load(): Promise<void> {
    const configs = await this.fetch();

    LANGUAGE_CODES.forEach((language) => {
      const { header, footer, learningEnvironment, services } = configs[language] ?? {};
      NAVIGATION_CONFIG[language] = {
        HEADER: header ?? [],
        FOOTER: footer ?? [],
        LEARNING_ENVIRONMENT: learningEnvironment ?? [],
        SERVICES: services ?? [],
      };
    });
  }
}

export default new NavigationService();
