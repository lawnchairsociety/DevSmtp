namespace DevSmtp.Core.Commands
{
    public sealed class SamlResult : CommandResult
    {
        public SamlResult()
        {
        }

        public SamlResult(Exception error)
            : base(error)
        {
        }
    }
}
