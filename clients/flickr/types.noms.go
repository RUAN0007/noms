// This file was generated by nomdl/codegen.

package main

import (
	"github.com/attic-labs/noms/chunks"
	"github.com/attic-labs/noms/ref"
	"github.com/attic-labs/noms/types"
)

var __mainPackageInFile_types_CachedRef ref.Ref

// This function builds up a Noms value that describes the type
// package implemented by this file and registers it with the global
// type package definition cache.
func init() {
	p := types.NewPackage([]types.Type{
		types.MakeStructType("User",
			[]types.Field{
				types.Field{"Id", types.MakePrimitiveType(types.StringKind), false},
				types.Field{"Name", types.MakePrimitiveType(types.StringKind), false},
				types.Field{"OAuthToken", types.MakePrimitiveType(types.StringKind), false},
				types.Field{"OAuthSecret", types.MakePrimitiveType(types.StringKind), false},
				types.Field{"Albums", types.MakeCompoundType(types.MapKind, types.MakePrimitiveType(types.StringKind), types.MakeType(ref.Ref{}, 1)), false},
			},
			types.Choices{},
		),
		types.MakeStructType("Album",
			[]types.Field{
				types.Field{"Id", types.MakePrimitiveType(types.StringKind), false},
				types.Field{"Title", types.MakePrimitiveType(types.StringKind), false},
				types.Field{"Photos", types.MakeCompoundType(types.RefKind, types.MakeCompoundType(types.SetKind, types.MakeCompoundType(types.RefKind, types.MakeType(ref.Parse("sha1-42722c55cec055ff32526d1dc0636fa4d48c9fd5"), 0)))), false},
			},
			types.Choices{},
		),
	}, []ref.Ref{
		ref.Parse("sha1-42722c55cec055ff32526d1dc0636fa4d48c9fd5"),
	})
	__mainPackageInFile_types_CachedRef = types.RegisterPackage(&p)
}

// User

type User struct {
	_Id          string
	_Name        string
	_OAuthToken  string
	_OAuthSecret string
	_Albums      MapOfStringToAlbum

	cs  chunks.ChunkStore
	ref *ref.Ref
}

func NewUser(cs chunks.ChunkStore) User {
	return User{
		_Id:          "",
		_Name:        "",
		_OAuthToken:  "",
		_OAuthSecret: "",
		_Albums:      NewMapOfStringToAlbum(cs),

		cs:  cs,
		ref: &ref.Ref{},
	}
}

type UserDef struct {
	Id          string
	Name        string
	OAuthToken  string
	OAuthSecret string
	Albums      MapOfStringToAlbumDef
}

func (def UserDef) New(cs chunks.ChunkStore) User {
	return User{
		_Id:          def.Id,
		_Name:        def.Name,
		_OAuthToken:  def.OAuthToken,
		_OAuthSecret: def.OAuthSecret,
		_Albums:      def.Albums.New(cs),
		cs:           cs,
		ref:          &ref.Ref{},
	}
}

func (s User) Def() (d UserDef) {
	d.Id = s._Id
	d.Name = s._Name
	d.OAuthToken = s._OAuthToken
	d.OAuthSecret = s._OAuthSecret
	d.Albums = s._Albums.Def()
	return
}

var __typeForUser types.Type

func (m User) Type() types.Type {
	return __typeForUser
}

func init() {
	__typeForUser = types.MakeType(__mainPackageInFile_types_CachedRef, 0)
	types.RegisterStruct(__typeForUser, builderForUser, readerForUser)
}

func builderForUser(cs chunks.ChunkStore, values []types.Value) types.Value {
	i := 0
	s := User{ref: &ref.Ref{}, cs: cs}
	s._Id = values[i].(types.String).String()
	i++
	s._Name = values[i].(types.String).String()
	i++
	s._OAuthToken = values[i].(types.String).String()
	i++
	s._OAuthSecret = values[i].(types.String).String()
	i++
	s._Albums = values[i].(MapOfStringToAlbum)
	i++
	return s
}

func readerForUser(v types.Value) []types.Value {
	values := []types.Value{}
	s := v.(User)
	values = append(values, types.NewString(s._Id))
	values = append(values, types.NewString(s._Name))
	values = append(values, types.NewString(s._OAuthToken))
	values = append(values, types.NewString(s._OAuthSecret))
	values = append(values, s._Albums)
	return values
}

func (s User) Equals(other types.Value) bool {
	return other != nil && __typeForUser.Equals(other.Type()) && s.Ref() == other.Ref()
}

func (s User) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s User) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, __typeForUser.Chunks()...)
	chunks = append(chunks, s._Albums.Chunks()...)
	return
}

func (s User) ChildValues() (ret []types.Value) {
	ret = append(ret, types.NewString(s._Id))
	ret = append(ret, types.NewString(s._Name))
	ret = append(ret, types.NewString(s._OAuthToken))
	ret = append(ret, types.NewString(s._OAuthSecret))
	ret = append(ret, s._Albums)
	return
}

func (s User) Id() string {
	return s._Id
}

func (s User) SetId(val string) User {
	s._Id = val
	s.ref = &ref.Ref{}
	return s
}

func (s User) Name() string {
	return s._Name
}

func (s User) SetName(val string) User {
	s._Name = val
	s.ref = &ref.Ref{}
	return s
}

func (s User) OAuthToken() string {
	return s._OAuthToken
}

func (s User) SetOAuthToken(val string) User {
	s._OAuthToken = val
	s.ref = &ref.Ref{}
	return s
}

func (s User) OAuthSecret() string {
	return s._OAuthSecret
}

func (s User) SetOAuthSecret(val string) User {
	s._OAuthSecret = val
	s.ref = &ref.Ref{}
	return s
}

func (s User) Albums() MapOfStringToAlbum {
	return s._Albums
}

func (s User) SetAlbums(val MapOfStringToAlbum) User {
	s._Albums = val
	s.ref = &ref.Ref{}
	return s
}

// Album

type Album struct {
	_Id     string
	_Title  string
	_Photos RefOfSetOfRefOfRemotePhoto

	cs  chunks.ChunkStore
	ref *ref.Ref
}

func NewAlbum(cs chunks.ChunkStore) Album {
	return Album{
		_Id:     "",
		_Title:  "",
		_Photos: NewRefOfSetOfRefOfRemotePhoto(ref.Ref{}),

		cs:  cs,
		ref: &ref.Ref{},
	}
}

type AlbumDef struct {
	Id     string
	Title  string
	Photos ref.Ref
}

func (def AlbumDef) New(cs chunks.ChunkStore) Album {
	return Album{
		_Id:     def.Id,
		_Title:  def.Title,
		_Photos: NewRefOfSetOfRefOfRemotePhoto(def.Photos),
		cs:      cs,
		ref:     &ref.Ref{},
	}
}

func (s Album) Def() (d AlbumDef) {
	d.Id = s._Id
	d.Title = s._Title
	d.Photos = s._Photos.TargetRef()
	return
}

var __typeForAlbum types.Type

func (m Album) Type() types.Type {
	return __typeForAlbum
}

func init() {
	__typeForAlbum = types.MakeType(__mainPackageInFile_types_CachedRef, 1)
	types.RegisterStruct(__typeForAlbum, builderForAlbum, readerForAlbum)
}

func builderForAlbum(cs chunks.ChunkStore, values []types.Value) types.Value {
	i := 0
	s := Album{ref: &ref.Ref{}, cs: cs}
	s._Id = values[i].(types.String).String()
	i++
	s._Title = values[i].(types.String).String()
	i++
	s._Photos = values[i].(RefOfSetOfRefOfRemotePhoto)
	i++
	return s
}

func readerForAlbum(v types.Value) []types.Value {
	values := []types.Value{}
	s := v.(Album)
	values = append(values, types.NewString(s._Id))
	values = append(values, types.NewString(s._Title))
	values = append(values, s._Photos)
	return values
}

func (s Album) Equals(other types.Value) bool {
	return other != nil && __typeForAlbum.Equals(other.Type()) && s.Ref() == other.Ref()
}

func (s Album) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s Album) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, __typeForAlbum.Chunks()...)
	chunks = append(chunks, s._Photos.Chunks()...)
	return
}

func (s Album) ChildValues() (ret []types.Value) {
	ret = append(ret, types.NewString(s._Id))
	ret = append(ret, types.NewString(s._Title))
	ret = append(ret, s._Photos)
	return
}

func (s Album) Id() string {
	return s._Id
}

func (s Album) SetId(val string) Album {
	s._Id = val
	s.ref = &ref.Ref{}
	return s
}

func (s Album) Title() string {
	return s._Title
}

func (s Album) SetTitle(val string) Album {
	s._Title = val
	s.ref = &ref.Ref{}
	return s
}

func (s Album) Photos() RefOfSetOfRefOfRemotePhoto {
	return s._Photos
}

func (s Album) SetPhotos(val RefOfSetOfRefOfRemotePhoto) Album {
	s._Photos = val
	s.ref = &ref.Ref{}
	return s
}

// MapOfStringToAlbum

type MapOfStringToAlbum struct {
	m   types.Map
	cs  chunks.ChunkStore
	ref *ref.Ref
}

func NewMapOfStringToAlbum(cs chunks.ChunkStore) MapOfStringToAlbum {
	return MapOfStringToAlbum{types.NewTypedMap(cs, __typeForMapOfStringToAlbum), cs, &ref.Ref{}}
}

type MapOfStringToAlbumDef map[string]AlbumDef

func (def MapOfStringToAlbumDef) New(cs chunks.ChunkStore) MapOfStringToAlbum {
	kv := make([]types.Value, 0, len(def)*2)
	for k, v := range def {
		kv = append(kv, types.NewString(k), v.New(cs))
	}
	return MapOfStringToAlbum{types.NewTypedMap(cs, __typeForMapOfStringToAlbum, kv...), cs, &ref.Ref{}}
}

func (m MapOfStringToAlbum) Def() MapOfStringToAlbumDef {
	def := make(map[string]AlbumDef)
	m.m.Iter(func(k, v types.Value) bool {
		def[k.(types.String).String()] = v.(Album).Def()
		return false
	})
	return def
}

func (m MapOfStringToAlbum) Equals(other types.Value) bool {
	return other != nil && __typeForMapOfStringToAlbum.Equals(other.Type()) && m.Ref() == other.Ref()
}

func (m MapOfStringToAlbum) Ref() ref.Ref {
	return types.EnsureRef(m.ref, m)
}

func (m MapOfStringToAlbum) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, m.Type().Chunks()...)
	chunks = append(chunks, m.m.Chunks()...)
	return
}

func (m MapOfStringToAlbum) ChildValues() []types.Value {
	return append([]types.Value{}, m.m.ChildValues()...)
}

// A Noms Value that describes MapOfStringToAlbum.
var __typeForMapOfStringToAlbum types.Type

func (m MapOfStringToAlbum) Type() types.Type {
	return __typeForMapOfStringToAlbum
}

func init() {
	__typeForMapOfStringToAlbum = types.MakeCompoundType(types.MapKind, types.MakePrimitiveType(types.StringKind), types.MakeType(__mainPackageInFile_types_CachedRef, 1))
	types.RegisterValue(__typeForMapOfStringToAlbum, builderForMapOfStringToAlbum, readerForMapOfStringToAlbum)
}

func builderForMapOfStringToAlbum(cs chunks.ChunkStore, v types.Value) types.Value {
	return MapOfStringToAlbum{v.(types.Map), cs, &ref.Ref{}}
}

func readerForMapOfStringToAlbum(v types.Value) types.Value {
	return v.(MapOfStringToAlbum).m
}

func (m MapOfStringToAlbum) Empty() bool {
	return m.m.Empty()
}

func (m MapOfStringToAlbum) Len() uint64 {
	return m.m.Len()
}

func (m MapOfStringToAlbum) Has(p string) bool {
	return m.m.Has(types.NewString(p))
}

func (m MapOfStringToAlbum) Get(p string) Album {
	return m.m.Get(types.NewString(p)).(Album)
}

func (m MapOfStringToAlbum) MaybeGet(p string) (Album, bool) {
	v, ok := m.m.MaybeGet(types.NewString(p))
	if !ok {
		return NewAlbum(m.cs), false
	}
	return v.(Album), ok
}

func (m MapOfStringToAlbum) Set(k string, v Album) MapOfStringToAlbum {
	return MapOfStringToAlbum{m.m.Set(types.NewString(k), v), m.cs, &ref.Ref{}}
}

// TODO: Implement SetM?

func (m MapOfStringToAlbum) Remove(p string) MapOfStringToAlbum {
	return MapOfStringToAlbum{m.m.Remove(types.NewString(p)), m.cs, &ref.Ref{}}
}

type MapOfStringToAlbumIterCallback func(k string, v Album) (stop bool)

func (m MapOfStringToAlbum) Iter(cb MapOfStringToAlbumIterCallback) {
	m.m.Iter(func(k, v types.Value) bool {
		return cb(k.(types.String).String(), v.(Album))
	})
}

type MapOfStringToAlbumIterAllCallback func(k string, v Album)

func (m MapOfStringToAlbum) IterAll(cb MapOfStringToAlbumIterAllCallback) {
	m.m.IterAll(func(k, v types.Value) {
		cb(k.(types.String).String(), v.(Album))
	})
}

func (m MapOfStringToAlbum) IterAllP(concurrency int, cb MapOfStringToAlbumIterAllCallback) {
	m.m.IterAllP(concurrency, func(k, v types.Value) {
		cb(k.(types.String).String(), v.(Album))
	})
}

type MapOfStringToAlbumFilterCallback func(k string, v Album) (keep bool)

func (m MapOfStringToAlbum) Filter(cb MapOfStringToAlbumFilterCallback) MapOfStringToAlbum {
	out := m.m.Filter(func(k, v types.Value) bool {
		return cb(k.(types.String).String(), v.(Album))
	})
	return MapOfStringToAlbum{out, m.cs, &ref.Ref{}}
}

// RefOfUser

type RefOfUser struct {
	target ref.Ref
	ref    *ref.Ref
}

func NewRefOfUser(target ref.Ref) RefOfUser {
	return RefOfUser{target, &ref.Ref{}}
}

func (r RefOfUser) TargetRef() ref.Ref {
	return r.target
}

func (r RefOfUser) Ref() ref.Ref {
	return types.EnsureRef(r.ref, r)
}

func (r RefOfUser) Equals(other types.Value) bool {
	return other != nil && __typeForRefOfUser.Equals(other.Type()) && r.Ref() == other.Ref()
}

func (r RefOfUser) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, r.Type().Chunks()...)
	chunks = append(chunks, r.target)
	return
}

func (r RefOfUser) ChildValues() []types.Value {
	return nil
}

// A Noms Value that describes RefOfUser.
var __typeForRefOfUser types.Type

func (r RefOfUser) Type() types.Type {
	return __typeForRefOfUser
}

func (r RefOfUser) Less(other types.OrderedValue) bool {
	return r.TargetRef().Less(other.(types.RefBase).TargetRef())
}

func init() {
	__typeForRefOfUser = types.MakeCompoundType(types.RefKind, types.MakeType(__mainPackageInFile_types_CachedRef, 0))
	types.RegisterRef(__typeForRefOfUser, builderForRefOfUser)
}

func builderForRefOfUser(r ref.Ref) types.Value {
	return NewRefOfUser(r)
}

func (r RefOfUser) TargetValue(cs chunks.ChunkStore) User {
	return types.ReadValue(r.target, cs).(User)
}

func (r RefOfUser) SetTargetValue(val User, cs chunks.ChunkSink) RefOfUser {
	return NewRefOfUser(types.WriteValue(val, cs))
}

// RefOfSetOfRefOfRemotePhoto

type RefOfSetOfRefOfRemotePhoto struct {
	target ref.Ref
	ref    *ref.Ref
}

func NewRefOfSetOfRefOfRemotePhoto(target ref.Ref) RefOfSetOfRefOfRemotePhoto {
	return RefOfSetOfRefOfRemotePhoto{target, &ref.Ref{}}
}

func (r RefOfSetOfRefOfRemotePhoto) TargetRef() ref.Ref {
	return r.target
}

func (r RefOfSetOfRefOfRemotePhoto) Ref() ref.Ref {
	return types.EnsureRef(r.ref, r)
}

func (r RefOfSetOfRefOfRemotePhoto) Equals(other types.Value) bool {
	return other != nil && __typeForRefOfSetOfRefOfRemotePhoto.Equals(other.Type()) && r.Ref() == other.Ref()
}

func (r RefOfSetOfRefOfRemotePhoto) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, r.Type().Chunks()...)
	chunks = append(chunks, r.target)
	return
}

func (r RefOfSetOfRefOfRemotePhoto) ChildValues() []types.Value {
	return nil
}

// A Noms Value that describes RefOfSetOfRefOfRemotePhoto.
var __typeForRefOfSetOfRefOfRemotePhoto types.Type

func (r RefOfSetOfRefOfRemotePhoto) Type() types.Type {
	return __typeForRefOfSetOfRefOfRemotePhoto
}

func (r RefOfSetOfRefOfRemotePhoto) Less(other types.OrderedValue) bool {
	return r.TargetRef().Less(other.(types.RefBase).TargetRef())
}

func init() {
	__typeForRefOfSetOfRefOfRemotePhoto = types.MakeCompoundType(types.RefKind, types.MakeCompoundType(types.SetKind, types.MakeCompoundType(types.RefKind, types.MakeType(ref.Parse("sha1-42722c55cec055ff32526d1dc0636fa4d48c9fd5"), 0))))
	types.RegisterRef(__typeForRefOfSetOfRefOfRemotePhoto, builderForRefOfSetOfRefOfRemotePhoto)
}

func builderForRefOfSetOfRefOfRemotePhoto(r ref.Ref) types.Value {
	return NewRefOfSetOfRefOfRemotePhoto(r)
}

func (r RefOfSetOfRefOfRemotePhoto) TargetValue(cs chunks.ChunkStore) SetOfRefOfRemotePhoto {
	return types.ReadValue(r.target, cs).(SetOfRefOfRemotePhoto)
}

func (r RefOfSetOfRefOfRemotePhoto) SetTargetValue(val SetOfRefOfRemotePhoto, cs chunks.ChunkSink) RefOfSetOfRefOfRemotePhoto {
	return NewRefOfSetOfRefOfRemotePhoto(types.WriteValue(val, cs))
}

// SetOfRefOfRemotePhoto

type SetOfRefOfRemotePhoto struct {
	s   types.Set
	cs  chunks.ChunkStore
	ref *ref.Ref
}

func NewSetOfRefOfRemotePhoto(cs chunks.ChunkStore) SetOfRefOfRemotePhoto {
	return SetOfRefOfRemotePhoto{types.NewTypedSet(cs, __typeForSetOfRefOfRemotePhoto), cs, &ref.Ref{}}
}

type SetOfRefOfRemotePhotoDef map[ref.Ref]bool

func (def SetOfRefOfRemotePhotoDef) New(cs chunks.ChunkStore) SetOfRefOfRemotePhoto {
	l := make([]types.Value, len(def))
	i := 0
	for d, _ := range def {
		l[i] = NewRefOfRemotePhoto(d)
		i++
	}
	return SetOfRefOfRemotePhoto{types.NewTypedSet(cs, __typeForSetOfRefOfRemotePhoto, l...), cs, &ref.Ref{}}
}

func (s SetOfRefOfRemotePhoto) Def() SetOfRefOfRemotePhotoDef {
	def := make(map[ref.Ref]bool, s.Len())
	s.s.Iter(func(v types.Value) bool {
		def[v.(RefOfRemotePhoto).TargetRef()] = true
		return false
	})
	return def
}

func (s SetOfRefOfRemotePhoto) Equals(other types.Value) bool {
	return other != nil && __typeForSetOfRefOfRemotePhoto.Equals(other.Type()) && s.Ref() == other.Ref()
}

func (s SetOfRefOfRemotePhoto) Ref() ref.Ref {
	return types.EnsureRef(s.ref, s)
}

func (s SetOfRefOfRemotePhoto) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, s.Type().Chunks()...)
	chunks = append(chunks, s.s.Chunks()...)
	return
}

func (s SetOfRefOfRemotePhoto) ChildValues() []types.Value {
	return append([]types.Value{}, s.s.ChildValues()...)
}

// A Noms Value that describes SetOfRefOfRemotePhoto.
var __typeForSetOfRefOfRemotePhoto types.Type

func (m SetOfRefOfRemotePhoto) Type() types.Type {
	return __typeForSetOfRefOfRemotePhoto
}

func init() {
	__typeForSetOfRefOfRemotePhoto = types.MakeCompoundType(types.SetKind, types.MakeCompoundType(types.RefKind, types.MakeType(ref.Parse("sha1-42722c55cec055ff32526d1dc0636fa4d48c9fd5"), 0)))
	types.RegisterValue(__typeForSetOfRefOfRemotePhoto, builderForSetOfRefOfRemotePhoto, readerForSetOfRefOfRemotePhoto)
}

func builderForSetOfRefOfRemotePhoto(cs chunks.ChunkStore, v types.Value) types.Value {
	return SetOfRefOfRemotePhoto{v.(types.Set), cs, &ref.Ref{}}
}

func readerForSetOfRefOfRemotePhoto(v types.Value) types.Value {
	return v.(SetOfRefOfRemotePhoto).s
}

func (s SetOfRefOfRemotePhoto) Empty() bool {
	return s.s.Empty()
}

func (s SetOfRefOfRemotePhoto) Len() uint64 {
	return s.s.Len()
}

func (s SetOfRefOfRemotePhoto) Has(p RefOfRemotePhoto) bool {
	return s.s.Has(p)
}

type SetOfRefOfRemotePhotoIterCallback func(p RefOfRemotePhoto) (stop bool)

func (s SetOfRefOfRemotePhoto) Iter(cb SetOfRefOfRemotePhotoIterCallback) {
	s.s.Iter(func(v types.Value) bool {
		return cb(v.(RefOfRemotePhoto))
	})
}

type SetOfRefOfRemotePhotoIterAllCallback func(p RefOfRemotePhoto)

func (s SetOfRefOfRemotePhoto) IterAll(cb SetOfRefOfRemotePhotoIterAllCallback) {
	s.s.IterAll(func(v types.Value) {
		cb(v.(RefOfRemotePhoto))
	})
}

func (s SetOfRefOfRemotePhoto) IterAllP(concurrency int, cb SetOfRefOfRemotePhotoIterAllCallback) {
	s.s.IterAllP(concurrency, func(v types.Value) {
		cb(v.(RefOfRemotePhoto))
	})
}

type SetOfRefOfRemotePhotoFilterCallback func(p RefOfRemotePhoto) (keep bool)

func (s SetOfRefOfRemotePhoto) Filter(cb SetOfRefOfRemotePhotoFilterCallback) SetOfRefOfRemotePhoto {
	out := s.s.Filter(func(v types.Value) bool {
		return cb(v.(RefOfRemotePhoto))
	})
	return SetOfRefOfRemotePhoto{out, s.cs, &ref.Ref{}}
}

func (s SetOfRefOfRemotePhoto) Insert(p ...RefOfRemotePhoto) SetOfRefOfRemotePhoto {
	return SetOfRefOfRemotePhoto{s.s.Insert(s.fromElemSlice(p)...), s.cs, &ref.Ref{}}
}

func (s SetOfRefOfRemotePhoto) Remove(p ...RefOfRemotePhoto) SetOfRefOfRemotePhoto {
	return SetOfRefOfRemotePhoto{s.s.Remove(s.fromElemSlice(p)...), s.cs, &ref.Ref{}}
}

func (s SetOfRefOfRemotePhoto) Union(others ...SetOfRefOfRemotePhoto) SetOfRefOfRemotePhoto {
	return SetOfRefOfRemotePhoto{s.s.Union(s.fromStructSlice(others)...), s.cs, &ref.Ref{}}
}

func (s SetOfRefOfRemotePhoto) First() RefOfRemotePhoto {
	return s.s.First().(RefOfRemotePhoto)
}

func (s SetOfRefOfRemotePhoto) fromStructSlice(p []SetOfRefOfRemotePhoto) []types.Set {
	r := make([]types.Set, len(p))
	for i, v := range p {
		r[i] = v.s
	}
	return r
}

func (s SetOfRefOfRemotePhoto) fromElemSlice(p []RefOfRemotePhoto) []types.Value {
	r := make([]types.Value, len(p))
	for i, v := range p {
		r[i] = v
	}
	return r
}

// RefOfRemotePhoto

type RefOfRemotePhoto struct {
	target ref.Ref
	ref    *ref.Ref
}

func NewRefOfRemotePhoto(target ref.Ref) RefOfRemotePhoto {
	return RefOfRemotePhoto{target, &ref.Ref{}}
}

func (r RefOfRemotePhoto) TargetRef() ref.Ref {
	return r.target
}

func (r RefOfRemotePhoto) Ref() ref.Ref {
	return types.EnsureRef(r.ref, r)
}

func (r RefOfRemotePhoto) Equals(other types.Value) bool {
	return other != nil && __typeForRefOfRemotePhoto.Equals(other.Type()) && r.Ref() == other.Ref()
}

func (r RefOfRemotePhoto) Chunks() (chunks []ref.Ref) {
	chunks = append(chunks, r.Type().Chunks()...)
	chunks = append(chunks, r.target)
	return
}

func (r RefOfRemotePhoto) ChildValues() []types.Value {
	return nil
}

// A Noms Value that describes RefOfRemotePhoto.
var __typeForRefOfRemotePhoto types.Type

func (r RefOfRemotePhoto) Type() types.Type {
	return __typeForRefOfRemotePhoto
}

func (r RefOfRemotePhoto) Less(other types.OrderedValue) bool {
	return r.TargetRef().Less(other.(types.RefBase).TargetRef())
}

func init() {
	__typeForRefOfRemotePhoto = types.MakeCompoundType(types.RefKind, types.MakeType(ref.Parse("sha1-42722c55cec055ff32526d1dc0636fa4d48c9fd5"), 0))
	types.RegisterRef(__typeForRefOfRemotePhoto, builderForRefOfRemotePhoto)
}

func builderForRefOfRemotePhoto(r ref.Ref) types.Value {
	return NewRefOfRemotePhoto(r)
}

func (r RefOfRemotePhoto) TargetValue(cs chunks.ChunkStore) RemotePhoto {
	return types.ReadValue(r.target, cs).(RemotePhoto)
}

func (r RefOfRemotePhoto) SetTargetValue(val RemotePhoto, cs chunks.ChunkSink) RefOfRemotePhoto {
	return NewRefOfRemotePhoto(types.WriteValue(val, cs))
}
