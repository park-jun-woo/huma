package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNodeAnalyzer_Express(t *testing.T) {
	src := `const express = require('express');

function createUser(req, res) {
    if (!req.body.name) {
        return res.status(400).json({ error: 'name required' });
    }
    res.status(201).json(user);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.js")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "createUser", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[400] {
		t.Fatal("expected 400")
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
}

func TestNodeAnalyzer_SendStatus(t *testing.T) {
	src := `function deleteUser(req, res) {
    res.sendStatus(204);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.js")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "deleteUser", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 204 {
		t.Fatalf("expected 204, got %d", branches[0].Status)
	}
}

func TestNodeAnalyzer_Fastify(t *testing.T) {
	src := `async function createItem(request, reply) {
    if (!request.body.name) {
        return reply.code(400).send({ error: 'invalid' });
    }
    reply.code(201).send(item);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.js")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "createItem", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[400] {
		t.Fatal("expected 400")
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
}

func TestNodeAnalyzer_NestJSDecorator(t *testing.T) {
	src := `import { HttpCode } from '@nestjs/common';

@HttpCode(201)
async create(@Body() dto: CreateDto) {
    return this.service.create(dto);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "controller.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "create", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 201 {
		t.Fatalf("expected 201, got %d", branches[0].Status)
	}
}

func TestNodeAnalyzer_NestJSExceptions(t *testing.T) {
	src := `import { NotFoundException, BadRequestException } from '@nestjs/common';

async findOne(id: string) {
    if (!id) {
        throw new BadRequestException('id required');
    }
    const item = await this.repo.findOne(id);
    if (!item) {
        throw new NotFoundException('item not found');
    }
    return item;
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "service.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "findOne", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[400] {
		t.Fatal("expected 400")
	}
	if !statuses[404] {
		t.Fatal("expected 404")
	}
}

func TestNodeAnalyzer_HttpException(t *testing.T) {
	src := `import { HttpException } from '@nestjs/common';

async update(id: string, dto: UpdateDto) {
    throw new HttpException('forbidden', 403);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "service.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "update", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 403 {
		t.Fatalf("expected 403, got %d", branches[0].Status)
	}
}

func TestNodeAnalyzer_LineRange(t *testing.T) {
	src := `line1
    res.status(200).json(data);
line3
    res.status(404).json({ error: 'not found' });
line5
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.js")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "handler", 1, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch (line range filter), got %d", len(branches))
	}
	if branches[0].Status != 200 {
		t.Fatalf("expected 200, got %d", branches[0].Status)
	}
}

func TestNodeAnalyzer_InvalidFile(t *testing.T) {
	a := &NodeAnalyzer{}
	_, err := a.Analyze("/nonexistent/file.js", "handler", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestNodeAnalyzer_AllNestExceptions(t *testing.T) {
	src := `async handler() {
    throw new UnauthorizedException('no auth');
    throw new ForbiddenException('forbidden');
    throw new ConflictException('conflict');
    throw new InternalServerErrorException('error');
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "service.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 4 {
		t.Fatalf("expected 4 branches, got %d", len(branches))
	}

	expected := map[int]bool{401: true, 403: true, 409: true, 500: true}
	for _, b := range branches {
		if !expected[b.Status] {
			t.Fatalf("unexpected status %d", b.Status)
		}
	}
}

func TestNodeAnalyzer_ExtendedNestExceptions(t *testing.T) {
	src := `async handler() {
    throw new MethodNotAllowedException('method not allowed');
    throw new NotAcceptableException('not acceptable');
    throw new RequestTimeoutException('timeout');
    throw new GoneException('gone');
    throw new PreconditionFailedException('precondition failed');
    throw new PayloadTooLargeException('too large');
    throw new UnsupportedMediaTypeException('unsupported');
    throw new UnprocessableEntityException('unprocessable');
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "service.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 8 {
		t.Fatalf("expected 8 branches, got %d", len(branches))
	}

	expected := map[int]bool{
		405: true, 406: true, 408: true, 410: true,
		412: true, 413: true, 415: true, 422: true,
	}
	for _, b := range branches {
		if !expected[b.Status] {
			t.Fatalf("unexpected status %d", b.Status)
		}
	}
}

func TestNodeAnalyzer_HttpStatusEnum(t *testing.T) {
	src := `async create(dto: CreateDto) {
    return res.status(HttpStatus.CREATED).json(result);
}

async find(id: string) {
    if (!item) {
        throw new HttpException('not found', HttpStatus.NOT_FOUND);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "controller.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "create", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[201] {
		t.Fatal("expected 201 from HttpStatus.CREATED")
	}
	if !statuses[404] {
		t.Fatal("expected 404 from HttpStatus.NOT_FOUND")
	}
}

func TestNodeAnalyzer_ApiResponse(t *testing.T) {
	src := `@ApiResponse({ status: 201, description: 'Created' })
@ApiResponse({ status: 400, description: 'Bad Request' })
async create(@Body() dto: CreateDto) {
    return this.service.create(dto);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "controller.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "create", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
	if !statuses[400] {
		t.Fatal("expected 400")
	}
}

func TestNodeAnalyzer_ApiShorthand(t *testing.T) {
	src := `@ApiOkResponse({ description: 'Success' })
@ApiCreatedResponse({ description: 'Created' })
@ApiNotFoundResponse({ description: 'Not Found' })
async handler() {
    return this.service.handle();
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "controller.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[200] {
		t.Fatal("expected 200 from @ApiOkResponse")
	}
	if !statuses[201] {
		t.Fatal("expected 201 from @ApiCreatedResponse")
	}
	if !statuses[404] {
		t.Fatal("expected 404 from @ApiNotFoundResponse")
	}
}

func TestNodeAnalyzer_ExpressImplicit200(t *testing.T) {
	src := `function getUsers(req, res) {
    const users = db.findAll();
    res.json(users);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.js")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "getUsers", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 200 {
		t.Fatalf("expected 200, got %d", branches[0].Status)
	}
}

func TestNodeAnalyzer_ExpressStatusJsonNoConflict(t *testing.T) {
	src := `function createUser(req, res) {
    res.status(201).json(user);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.js")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "createUser", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 201 {
		t.Fatalf("expected 201, got %d", branches[0].Status)
	}
}

func TestNodeAnalyzer_ExpressRedirect(t *testing.T) {
	src := `function handler(req, res) {
    res.redirect(301, '/new-url');
    res.redirect('/other-url');
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.js")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[301] {
		t.Fatal("expected 301")
	}
	if !statuses[302] {
		t.Fatal("expected 302")
	}
}

func TestNodeAnalyzer_CreateError(t *testing.T) {
	src := `const createError = require('http-errors');

function getUser(req, res, next) {
    const user = db.find(req.params.id);
    if (!user) {
        return next(createError(404, 'User not found'));
    }
    if (!user.active) {
        return next(createError(403, 'User inactive'));
    }
    res.json(user);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.js")
	os.WriteFile(file, []byte(src), 0o644)

	a := &NodeAnalyzer{}
	branches, err := a.Analyze(file, "getUser", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 3 {
		t.Fatalf("expected 3 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[404] {
		t.Fatal("expected 404")
	}
	if !statuses[403] {
		t.Fatal("expected 403")
	}
	if !statuses[200] {
		t.Fatal("expected 200 from res.json()")
	}
}
