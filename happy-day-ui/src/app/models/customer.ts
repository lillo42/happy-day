export interface Customer {
  id: string;
  name: string;
  comment: string;
  phones: Phone[];
  createdAt: Date;
  modifiedAt: Date;
}

export interface Phone {
  number: string;
}

export enum CustomerSort {
  IdAsc = "id_asc",
  IdDesc = "id_desc",
  NameAsc = "name_asc",
  NameDesc = "name_asc",
  CommentAsc = "comment_desc",
  CommentDesc = "comment_desc",
}
