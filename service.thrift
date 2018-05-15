namespace * ogcache

struct OpenGraph {
  1:string title,
  2:string type,
  3:string image,
  4:string url,
  5:string audio,
  6:string description,
  7:string determiner,
  8:string locale,
  9:list<string> locale_alternate,
  10:string site_name,
  11:string video
}

service OgService {
  OpenGraph inquiry(1:string url)
}