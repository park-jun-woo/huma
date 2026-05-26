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
