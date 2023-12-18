import { CONTENT_URL } from 'config';
import { get } from 'utils/fetch-client';

export type TocTree = Array<{
  value?: [string, string];
  children: TocTree;
}>;

export interface ContentUnitHeading {
  id: string;
  level: string;
  text: string;
}

export interface ContentUnitText {
  text: string;
}

export interface ContentUnitInfobox {
  text: string;
  style: 'info';
}

export type TableData = string[][];

export interface ContentUnitTable {
  hasheadercolumn: boolean;
  hasheaderrow: boolean;
  table: TableData;
}

export interface ContentUnitImage {
  alt: string;
  caption: string;
  link: string;
  src: string;
  width?: string;
  height?: string;
  min_width?: string;
  min_height?: string;
}

export interface ContentUnitList {
  text: string;
}

export interface ContentUnitMarkdown {
  markdown: string;
  text: string;
}

export type ContentUnitLine = any[];

export interface ContentUnitCode {
  code: string;
}

export interface ContentUnitEntityTypeHeadline {
  id: string;
  entitytype: string;
  link?: string;
}

export interface ContentUnitFieldDescription {
  id: string;
  field: string;
  entitytype: string;
}

interface EntityTypeAndFields {
  name?: string;
  comment?: string;
  fieldsByImportance: { [key: string]: string[] };
}

export interface ContentUnitCompletenessDocumentationTable {
  label: string;
  description: string;
  entityTypes: EntityTypeAndFields[];
}

type ContentUnit =
  | ContentUnitHeading
  | ContentUnitText
  | ContentUnitInfobox
  | ContentUnitTable
  | ContentUnitImage
  | ContentUnitList
  | ContentUnitMarkdown
  | ContentUnitLine
  | ContentUnitCode
  | ContentUnitEntityTypeHeadline
  | ContentUnitFieldDescription
  | ContentUnitCompletenessDocumentationTable;

export enum ContentBlockType {
  heading = 'heading',
  text = 'text',
  infobox = 'infobox',
  table = 'table',
  image = 'image',
  list = 'list',
  markdown = 'markdown',
  line = 'line',
  code = 'code',
  entityTypeHeadline = 'entityTypeHeadline',
  fieldDescription = 'fieldDescription',
  completenessDocumentationTable = 'completenessDocumentationTable',
}

export interface ContentBlock {
  type: ContentBlockType;
  content: ContentUnit;
}

export interface Content {
  title: string;
  alignment: string;
  flags: string[];
  toc: TocTree;
  blocks: ContentBlock[];
}

export interface ContentResponse {
  [language: string]: Content;
}

export class ContentService {
  async fetch(pageId: string): Promise<ContentResponse> {
    const [response] = await get<ContentResponse>({
      url: `${CONTENT_URL}/${pageId}`,
    });
    return response;
  }
}

export default new ContentService();
