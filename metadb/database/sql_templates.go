package database

const userInfoWriteTemplate = `INSERT OR REPLACE INTO user_metadata (
    artist_pastel_id,
    real_name,
    facebook_link,
    twitter_link,
    native_currency,
    location,
    primary_language,
    categories,
    biography,
    timestamp,
    signature,
    previous_block_hash,
    user_data_hash,
    avatar_image,
    avatar_filename,
    cover_photo_image,
    cover_photo_filename
) VALUES (
    '{{.ArtistPastelID}}',
    '{{.RealName}}',
    '{{.FacebookLink}}',
    '{{.TwitterLink}}',
    '{{.NativeCurrency}}',
    '{{.Location}}',
    '{{.PrimaryLanguage}}',
    '{{.Categories}}',
    '{{.Biography}}',
    {{.Timestamp}},
    '{{.Signature}}',
    '{{.PreviousBlockHash}}',
    '{{.UserdataHash}}',
    {{ if (eq .AvatarImage "")}}NULL,{{ else }}x'{{.AvatarImage}}',{{ end }}
    '{{.AvatarFilename}}',
    {{ if (eq .CoverPhotoImage "")}}NULL,{{ else }}x'{{.CoverPhotoImage}}',{{ end }}
    '{{.CoverPhotoFilename}}'
);`

const userInfoQueryTemplate = `SELECT * FROM user_metadata WHERE artist_pastel_id = '{{.}}';`
