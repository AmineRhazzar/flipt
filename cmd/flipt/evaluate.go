package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"go.flipt.io/flipt/rpc/flipt"
)

type evaluateCommand struct {
	address   string
	token     string
	namespace string
	flag      string
	entityId  string
	context   map[string]string
	verbose   bool
}

func newEvaluateCommand() *cobra.Command {
	evaluate := &evaluateCommand{}

	cmd := &cobra.Command{
		Use:   "evaluate",
		Short: "evaluate a flag",
		RunE:  evaluate.run,
	}

	cmd.Flags().StringVarP(
		&evaluate.address,
		"address", "a",
		"",
		"address of Flipt instance to evaluate against",
	)

	cmd.Flags().StringVarP(
		&evaluate.token,
		"token", "t",
		"",
		"client token used to authenticate access to remote Flipt instance when evaluating.",
	)

	cmd.Flags().StringVarP(
		&evaluate.namespace,
		"namespace", "n",
		"default",
		"namespace key of the flag",
	)

	cmd.Flags().StringVarP(
		&evaluate.flag,
		"flag", "f",
		"",
		"flag key",
	)

	cmd.Flags().StringVarP(
		&evaluate.entityId,
		"entity-id", "e",
		"",
		"Entity Id to run the evaluation for",
	)

	cmd.Flags().StringToStringVarP(
		&evaluate.context,
		"context", "c",
		map[string]string{},
		"evaluation context",
	)

	cmd.MarkFlagRequired("flag")
	cmd.MarkFlagRequired("address")

	return cmd
}

func (e *evaluateCommand) run(cmd *cobra.Command, args []string) error {
	if e.address == "" {
		return fmt.Errorf("specified address is empty")
	}

	client, err := fliptClient(e.address, e.token)
	if err != nil {
		return err
	}
	
	if(e.entityId == "") {
		e.entityId = uuid.NewString()
	}

	response, evalErr := client.Evaluate(cmd.Context(), &flipt.EvaluationRequest{
		FlagKey:      e.flag,
		NamespaceKey: e.namespace,
		EntityId:     e.entityId,
		Context:      e.context,
	})

	if evalErr != nil {
		return evalErr
	}

	jsonResponse, _ := json.MarshalIndent(response, "", "  ")
	
	fmt.Printf("%s\n", jsonResponse)
	return nil
}
