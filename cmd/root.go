/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "odebo",
	Short: "list the untagged EC2 instances",
	Long:  `Command will list the untagged resources, which will help us to investigate find the correct owners.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		getUntaggedInstances(region)
	},
}

var region string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.odebo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&region, "region", "r", "us-east-1", "AWS region to check")
}

// getUntaggedInstances fetches EC2 instances without tags
func getUntaggedInstances(region string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	client := ec2.NewFromConfig(cfg)

	// Describe instances
	input := &ec2.DescribeInstancesInput{}
	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to describe instances: %v", err)
	}

	untaggedInstances := []string{}

	// Iterate over instances
	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if len(instance.Tags) == 0 {
				untaggedInstances = append(untaggedInstances, *instance.InstanceId)
			}
		}
	}

	// Print results
	if len(untaggedInstances) > 0 {
		fmt.Println("Untagged EC2 instances:")
		for _, id := range untaggedInstances {
			fmt.Println("- ", id)
		}
	} else {
		fmt.Println("All EC2 instances are tagged in this region.")
	}
}
