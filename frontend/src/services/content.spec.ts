jest.mock('stencil-router-v2');

import { CONTENT_URL } from 'config';
import * as fetchClient from 'utils/fetch-client';
import contentService, { Content, ContentBlockType } from './content';

const requestSpy = jest.spyOn(fetchClient, 'request');
requestSpy.mockImplementation(async () => [null, new Headers()] as [any, Headers]);

const EXAMPLE_CONTENT: Content = {
  title: 'Impressum',
  alignment: 'center',
  flags: ['LEGAL_PAGE', 'SKIP_AUTH'],
  toc: [
    {
      value: ['Label', 'http://localurl#anchor1'],
      children: [{ value: ['Label', 'http://localurl#anchor1.1'], children: [] }],
    },
  ],
  blocks: [
    {
      type: ContentBlockType.heading,
      content: {
        level: 'h1',
        text: 'Impressum',
      },
    },
    {
      type: ContentBlockType.heading,
      content: {
        level: 'h3',
        text: 'Kontakt',
      },
    },
    {
      type: ContentBlockType.text,
      content: {
        text: '<p>E-Mail: <a href="mailto:mex@example.com">mex@example.com</a></p>',
      },
    },
    {
      type: ContentBlockType.heading,
      content: {
        level: 'h3',
        text: 'Hosting',
      },
    },
    {
      type: ContentBlockType.text,
      content: {
        text: '<p>D4L data4life gGmbH</p>',
      },
    },
    {
      type: ContentBlockType.image,
      content: {
        src: '/media/pages/pages/imprint/19208145ce-1666860766/screenshot-2022-10-20-at-09.18.35.png',
        alt: 'stats about LGBTQIA+',
        caption: 'my image',
        link: 'https://www.zdf.de/kinder/logo/das-bedeutet-lgbtqia-100.html',
        width: 'auto',
        height: 'auto',
        min_width: 'auto',
        min_height: 'auto',
      },
    },
    {
      type: ContentBlockType.markdown,
      content: {
        text: `<h1>Markdown test</h1>
<ul>
<li>one</li>
<li>two</li>
<li>three</li>
</ul>
<hr />
<p><code>Hello World!</code></p>
<p><img src="image.jpg" alt="alt text" /></p>`,
        markdown: `# Markdown test

- one
- two
- three

---
\`Hello World!\`

![alt text](image.jpg)`,
      },
    },
    {
      type: ContentBlockType.list,
      content: {
        text: '<ul><li>foo</li><li>bar</li><li>baz</li><li></li></ul>',
      },
    },
    {
      type: ContentBlockType.line,
      content: [],
    },
    {
      type: ContentBlockType.code,
      content: {
        code: 'console.log("hello");',
      },
    },
  ],
};

describe('content service', () => {
  describe('fetch()', () => {
    it('fetches content data by page ID', async () => {
      requestSpy.mockImplementationOnce(async () => [EXAMPLE_CONTENT, new Headers()]);

      expect(await contentService.fetch('foo-id')).toBe(EXAMPLE_CONTENT);
      expect(requestSpy).toHaveBeenCalledWith({
        method: 'GET',
        url: `${CONTENT_URL}/foo-id`,
      });
    });

    it('throws an error when the request fails', async () => {
      const consoleSpy = jest.spyOn(console, 'error');
      consoleSpy.mockImplementationOnce(() => {});
      requestSpy.mockImplementationOnce(() => {
        throw new Error('Network error');
      });

      await expect(contentService.fetch('foo-id')).rejects.toThrow();
    });
  });
});
