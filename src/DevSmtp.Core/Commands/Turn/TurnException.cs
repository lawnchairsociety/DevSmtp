namespace DevSmtp.Core.Commands
{
    public class TurnException : Exception
    {
        public TurnException(string message)
            : base(message)
        {
        }

        public TurnException(string message, Exception innerException)
            : base(message, innerException)
        {
        }
    }
}
