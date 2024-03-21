export interface Resource {
  id: string;
  url: string;
  name: string;
  image_url?: null | string;
  author?: string | null;
  description?: string | null;
  status?: string | null;
  review_rating?: string | null;
  review_comment?: string | null;
  created_at: string;
}
